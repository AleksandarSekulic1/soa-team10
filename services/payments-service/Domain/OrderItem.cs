// services/payments-service/Domain/OrderItem.cs
namespace PaymentsService.Domain;

public class OrderItem
{
    public Guid Id { get; set; }
    public string TourName { get; set; }
    public double Price { get; set; }
    public string TourId { get; set; } 
}