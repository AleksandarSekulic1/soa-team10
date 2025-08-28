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

    // POST /api/shopping-cart/{touristUsername}/items
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

        // PROVERA DA LI TURA VEĆ POSTOJI U KORPI
        bool tourExists = cart.Items.Any(i => i.TourId == item.TourId);
        if (tourExists)
        {
            return Conflict(new { message = "Ova tura se već nalazi u korpi." });
        }

        cart.Items.Add(item);
        await _context.SaveChangesAsync();

        return Ok(cart);
    }

    // GET /api/shopping-cart/{touristUsername}
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

    // NOVA METODA: DELETE /api/shopping-cart/{touristUsername}/items/{itemId}
    [HttpDelete("{touristUsername}/items/{itemId:guid}")]
    public async Task<IActionResult> RemoveItemFromCart(string touristUsername, Guid itemId)
    {
        var cart = await _context.ShoppingCarts
            .Include(c => c.Items)
            .FirstOrDefaultAsync(c => c.TouristUsername == touristUsername);

        if (cart == null)
        {
            return NotFound(new { message = "Korpa nije pronađena." });
        }

        var itemToRemove = cart.Items.FirstOrDefault(i => i.Id == itemId);
        if (itemToRemove == null)
        {
            return NotFound(new { message = "Stavka nije pronađena u korpi." });
        }

        _context.OrderItems.Remove(itemToRemove);
        await _context.SaveChangesAsync();

        // Vraćamo ažuriranu korpu
        return Ok(cart);
    }
    [HttpPost("{touristUsername}/checkout")]
    public async Task<IActionResult> Checkout(string touristUsername)
    {
        var cart = await _context.ShoppingCarts
            .Include(c => c.Items)
            .FirstOrDefaultAsync(c => c.TouristUsername == touristUsername);

        if (cart == null || !cart.Items.Any())
        {
            return BadRequest(new { message = "Korpa je prazna." });
        }

        var purchaseTokens = new List<TourPurchaseToken>();
        foreach (var item in cart.Items)
        {
            purchaseTokens.Add(new TourPurchaseToken
            {
                TouristUsername = touristUsername, // ISPRAVLJENO: Koristimo username
                TourId = item.TourId,
                PurchaseTime = DateTime.UtcNow
            });
        }

        await _context.TourPurchaseTokens.AddRangeAsync(purchaseTokens);
        _context.OrderItems.RemoveRange(cart.Items);
        await _context.SaveChangesAsync();

        return Ok(new { message = "Kupovina uspešna!", tokens = purchaseTokens });
    }
    [HttpGet("{touristUsername}/tokens")]
    public async Task<IActionResult> GetPurchaseTokens(string touristUsername)
    {
        var tokens = await _context.TourPurchaseTokens
            .Where(t => t.TouristUsername == touristUsername)
            .ToListAsync();

        return Ok(tokens);
    }
}
