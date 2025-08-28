// services/payments-service/Controllers/ShoppingCartController.cs
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using PaymentsService.Data;
using PaymentsService.Domain;

namespace PaymentsService.Controllers;

[ApiController]
[Route("api/shopping-cart")]
public class ShoppingCartController : ControllerBase
{
    private readonly PaymentsDbContext _context;

    public ShoppingCartController(PaymentsDbContext context)
    {
        _context = context;
    }

    // Ruta sada prihvata username kao string
    [HttpPost("{touristUsername}/items")]
    public async Task<IActionResult> AddItemToCart(string touristUsername, [FromBody] OrderItem item)
    {
        var cart = await _context.ShoppingCarts
            .Include(c => c.Items)
            .FirstOrDefaultAsync(c => c.TouristUsername == touristUsername);

        if (cart == null)
        {
            cart = new ShoppingCart { TouristUsername = touristUsername };
            _context.ShoppingCarts.Add(cart);
        }

        cart.Items.Add(item);
        await _context.SaveChangesAsync();

        return Ok(cart);
    }

    // Ruta sada prihvata username kao string
    [HttpGet("{touristUsername}")]
    public async Task<IActionResult> GetCart(string touristUsername)
    {
        var cart = await _context.ShoppingCarts
            .Include(c => c.Items)
            .FirstOrDefaultAsync(c => c.TouristUsername == touristUsername);

        if (cart == null)
        {
            return Ok(new ShoppingCart { TouristUsername = touristUsername });
        }

        return Ok(cart);
    }
}