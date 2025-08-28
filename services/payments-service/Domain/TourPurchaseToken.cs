namespace PaymentsService.Domain;

public class TourPurchaseToken
{
    public Guid Id { get; set; }
    public string TouristUsername { get; set; } 
    public string TourId { get; set; }        
    public DateTime PurchaseTime { get; set; }
}
