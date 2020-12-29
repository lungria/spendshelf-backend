package category

const (
	// common categories

	// Default is the ID of category, that must be used for all new imported transactions.
	Default = 1
	// Ignored is the ID of category, that would be ignored in report.
	Ignored = 127

	// todo: remove hardcoded
	// custom categories, used in automatcher package (optional)

	// Health is for transactions, related to medicine, doctors appointments and care products.
	Health = 2
	// Grocery is for shops, groceries delivery.
	Grocery = 3
	// Tax is for government taxes.
	Tax = 4
	// Rent is for apartment related spendings.
	Rent = 5
	// Internet is for home and mobile internet related spendings.
	Internet = 6
	// Clothes is for clothes and footwear spendings.
	Clothes = 7
	// IKEA is a generic trademark for all furniture, housekeeping products and etc.
	IKEA = 8
	// Caffee is for prepared meal outside of home and food delivery services.
	Caffee = 9
	// Other transactions, that can't be placed in existing groups
	Other = 10
)
