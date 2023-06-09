package api

import (
	"github.com/betelgeusexru/golang-hotel-reservation/db"
	"github.com/betelgeusexru/golang-hotel-reservation/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID || !user.IsAdmin {
		return fiber.NewError(401, "not authorized")
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.Status(200).JSON(map[string]string{"msg": "updated"})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return fiber.NewError(401, "not authorized")
	}

	return c.JSON(booking)
}