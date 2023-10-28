package hostel

type HostelQueryParams struct {
	Name                 string   `form:"name" binding:"omitempty"`
	Description          string   `form:"description" binding:"omitempty"`
	UniversityID         uint     `form:"university_id" binding:"omitempty"`
	Address              string   `form:"address" binding:"omitempty"`
	City                 string   `form:"city" binding:"omitempty"`
	State                string   `form:"state" binding:"omitempty"`
	Country              string   `form:"country" binding:"omitempty"`
	ManagerID            uint     `form:"manager_id" binding:"omitempty"`
	NumberOfUnits        uint     `form:"number_of_units" binding:"omitempty"`
	NumberOfBedrooms     uint     `form:"number_of_bedrooms" binding:"omitempty"`
	NumberOfBathrooms    uint     `form:"number_of_bathrooms" binding:"omitempty"`
	Kitchen              string   `form:"kitchen" binding:"omitempty,oneof=shared none private"`
	FloorSpace           uint     `form:"floor_space" binding:"omitempty"`
	HostelFeeTotalMin    float64  `form:"hostel_fee_total_min" binding:"omitempty"`
	HostelFeeTotalMax    float64  `form:"hostel_fee_total_max" binding:"omitempty"`
	HostelFeePlan        string   `form:"hostel_fee_plan" binding:"omitempty"`
	HasAmenities         bool     `form:"has_amenities" binding:"omitempty"`
	SecurityRatingMin    *float32 `form:"security_rating_min" binding:"omitempty"`
	SecurityRatingMax    *float32 `form:"security_rating_max" binding:"omitempty"`
	LocationRatingMin    *float32 `form:"location_rating_min" binding:"omitempty"`
	LocationRatingMax    *float32 `form:"location_rating_max" binding:"omitempty"`
	GeneralRatingMin     *float32 `form:"general_rating_min" binding:"omitempty"`
	GeneralRatingMax     *float32 `form:"general_rating_max" binding:"omitempty"`
	AmenitiesRatingMin   *float32 `form:"amenities_rating_min" binding:"omitempty"`
	AmenitiesRatingMax   *float32 `form:"amenities_rating_max" binding:"omitempty"`
	WaterRatingMin       *float32 `form:"water_rating_min" binding:"omitempty"`
	WaterRatingMax       *float32 `form:"water_rating_max" binding:"omitempty"`
	ElectricityRatingMin *float32 `form:"electricity_rating_min" binding:"omitempty"`
	ElectricityRatingMax *float32 `form:"electricity_rating_max" binding:"omitempty"`
	CaretakerRatingMin   *float32 `form:"caretaker_rating_min" binding:"omitempty"`
	CaretakerRatingMax   *float32 `form:"caretaker_rating_max" binding:"omitempty"`
}
