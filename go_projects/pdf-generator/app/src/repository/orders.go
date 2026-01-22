package repository

type Item struct {
	Name     string
	Quantity int
	Price    float64
}

func (i Item) Total() float64 {
	return float64(i.Quantity) * i.Price
}

type Order struct {
	ID       int
	Customer string
	Date     string
	Items    []Item
}

var mockOrders = []Order{
	{
		ID:       101,
		Customer: "Анна Иванова",
		Date:     "2025-06-08",
		Items: []Item{
			{Name: "Ноутбук", Quantity: 1, Price: 74990.00},
			{Name: "Мышь", Quantity: 2, Price: 1290.00},
			{Name: "Коврик", Quantity: 1, Price: 590.00},
		},
	},
	{
		ID:       102,
		Customer: "Борис Смирнов",
		Date:     "2025-06-07",
		Items: []Item{
			{Name: "Монитор", Quantity: 2, Price: 18990.00},
		},
	},
	{
		ID:       103,
		Customer: "Вера Козлова",
		Date:     "2025-06-06",
		Items: []Item{
			{Name: "Клавиатура", Quantity: 1, Price: 3990.00},
			{Name: "Наушники", Quantity: 1, Price: 5490.00},
		},
	},
	{
		ID:       104,
		Customer: "Григорий Орлов",
		Date:     "2025-06-05",
		Items: []Item{
			{Name: "Смартфон", Quantity: 1, Price: 69990.00},
		},
	},
	{
		ID:       105,
		Customer: "Дарья Петрова",
		Date:     "2025-06-04",
		Items: []Item{
			{Name: "Планшет", Quantity: 1, Price: 45990.00},
			{Name: "Чехол", Quantity: 1, Price: 1990.00},
		},
	},
}

type OrdersRepository struct{}

func NewRepository() *OrdersRepository {
	return &OrdersRepository{}
}

func (r *OrdersRepository) GetOrderByID(id int) (Order, bool) {
	for _, o := range mockOrders {
		if o.ID == id {
			return o, true
		}
	}
	return Order{}, false
}
