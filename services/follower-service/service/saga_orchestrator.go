package service

import (
	"context"
	"fmt"
	"log"

	"follower-service/client"
	"follower-service/repository"
)

// SagaOrchestrator upravlja saga obrascem za unfollow operacije
type SagaOrchestrator struct {
	followerRepo *repository.FollowerRepository
	blogClient   *client.BlogClient
}

// SagaStep predstavlja jedan korak u saga transakciji
type SagaStep struct {
	Execute    func() error
	Compensate func() error
	Name       string
}

// UnfollowSaga konfiguracija za unfollow saga transakciju
type UnfollowSaga struct {
	FollowerID string
	FollowingID string
	Steps       []SagaStep
}

func NewSagaOrchestrator(followerRepo *repository.FollowerRepository, blogClient *client.BlogClient) *SagaOrchestrator {
	return &SagaOrchestrator{
		followerRepo: followerRepo,
		blogClient:   blogClient,
	}
}

// ExecuteUnfollowSaga izvršava saga obrazac za unfollow operaciju
func (s *SagaOrchestrator) ExecuteUnfollowSaga(ctx context.Context, followerID, followingID string) error {
	saga := &UnfollowSaga{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	// Korak 1: Uklanjanje follow veze iz baze
	saga.Steps = append(saga.Steps, SagaStep{
		Name: "Remove follow relationship",
		Execute: func() error {
			log.Printf("Saga Step 1: Removing follow relationship between %s and %s", followerID, followingID)
			return s.followerRepo.UnfollowUser(ctx, followerID, followingID)
		},
		Compensate: func() error {
			log.Printf("Saga Compensation 1: Restoring follow relationship between %s and %s", followerID, followingID)
			return s.followerRepo.FollowUser(ctx, followerID, followingID)
		},
	})

	// Korak 2: Uklanjanje lajkova iz blog servisa
	saga.Steps = append(saga.Steps, SagaStep{
		Name: "Remove likes from author's blogs",
		Execute: func() error {
			log.Printf("Saga Step 2: Removing likes from %s's blogs by user %s", followingID, followerID)
			return s.blogClient.RemoveLikesFromAuthorBlogs(followerID, followingID)
		},
		Compensate: func() error {
			log.Printf("Saga Compensation 2: Cannot restore likes - this operation is not reversible")
			return nil
		},
	})

	return s.executeSaga(saga)
}

// executeSaga izvršava sve korake saga transakcije
func (s *SagaOrchestrator) executeSaga(saga *UnfollowSaga) error {
	executedSteps := 0

	// Izvršavamo korake jedan po jedan
	for i, step := range saga.Steps {
		log.Printf("Executing saga step %d: %s", i+1, step.Name)
		
		if err := step.Execute(); err != nil {
			log.Printf("Saga step %d failed: %v", i+1, err)
			
			// Ako je korak neuspešan, izvršavamo kompenzacije za sve prethodno izvršene korake
			s.compensate(saga, executedSteps)
			return fmt.Errorf("saga failed at step %d (%s): %v", i+1, step.Name, err)
		}
		
		executedSteps++
		log.Printf("Saga step %d completed successfully", i+1)
	}

	log.Printf("Saga completed successfully for unfollow operation: %s -> %s", saga.FollowerID, saga.FollowingID)
	return nil
}

// compensate izvršava kompenzacije za neuspešnu saga transakciju
func (s *SagaOrchestrator) compensate(saga *UnfollowSaga, executedSteps int) {
	log.Printf("Starting saga compensation for %d executed steps", executedSteps)
	
	// Izvršavamo kompenzacije u obrnutom redosledu
	for i := executedSteps - 1; i >= 0; i-- {
		step := saga.Steps[i]
		log.Printf("Executing compensation for step %d: %s", i+1, step.Name)
		
		if err := step.Compensate(); err != nil {
			log.Printf("Compensation failed for step %d: %v", i+1, err)
			// U produkciji, ovo bi trebalo da se loguje u monitoring sistem
		} else {
			log.Printf("Compensation for step %d completed successfully", i+1)
		}
	}
	
	log.Printf("Saga compensation completed")
}