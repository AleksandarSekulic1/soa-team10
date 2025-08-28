// services/payments-service/Domain/ShoppingCart.cs
namespace PaymentsService.Domain;

public class ShoppingCart
{
    public Guid Id { get; set; }
    public string TouristUsername { get; set; } // <-- PROMENJENO
    public List<OrderItem> Items { get; set; } = new();
    public double TotalPrice => Items.Sum(i => i.Price);
}