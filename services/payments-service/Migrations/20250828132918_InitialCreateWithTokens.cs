using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace payments_service.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreateWithTokens : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "TourPurchaseTokens",
                columns: table => new
                {
                    Id = table.Column<Guid>(type: "uuid", nullable: false),
                    TouristUsername = table.Column<string>(type: "text", nullable: false),
                    TourId = table.Column<string>(type: "text", nullable: false),
                    PurchaseTime = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_TourPurchaseTokens", x => x.Id);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "TourPurchaseTokens");
        }
    }
}
