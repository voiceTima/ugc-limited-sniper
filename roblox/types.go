package roblox

type CurrentUserResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type CatalogDetailsResponse struct {
	Data []POSTItemDetails `json:"data"`
}

type POSTItemDetails struct {
	ID                           int      `json:"id"`
	ItemType                     string   `json:"itemType"`
	AssetType                    int      `json:"assetType"`
	Name                         string   `json:"name"`
	Description                  string   `json:"description"`
	Genres                       []string `json:"genres"`
	ItemRestrictions             []string `json:"itemRestrictions"`
	CreatorHasVerifiedBadge      bool     `json:"creatorHasVerifiedBadge"`
	CreatorType                  string   `json:"creatorType"`
	CreatorTargetId              int      `json:"creatorTargetId"`
	CreatorName                  string   `json:"creatorName"`
	Price                        int      `json:"price"`
	PriceStatus                  string   `json:"priceStatus"`
	UnitsAvailableForConsumption int      `json:"unitsAvailableForConsumption"`
	FavoriteCount                int      `json:"favoriteCount"`
	OffSaleDeadline              *string  `json:"offSaleDeadline"`
	CollectibleItemId            string   `json:"collectibleItemId"`
	TotalQuantity                int      `json:"totalQuantity"`
	SaleLocationType             string   `json:"saleLocationType"`
}

type PurchaseItemResponse struct {
	PurchaseResult string  `json:"purchaseResult"`
	Purchased      bool    `json:"purchased"`
	ErrorMessage   *string `json:"errorMessage"`
}
