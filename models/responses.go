package models

// ShortlinkResponseData structure
type ShortlinkResponseData struct {
	ID    uint64 `json:"id"`
	Short string `json:"short"`
	Full  string `json:"full"`
}

// UserResponseData contains information about user
type UserResponseData struct {
	ID   uint64 `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique;not null" binding:"required,min=5"`
}

// UserResponse contains information about user
type UserResponse struct {
	Data   UserResponseData `json:"data"`
	Result string           `json:"result"`
}

// FullLinkUseCountResponse structure
type FullLinkUseCountResponse struct {
	FullLink  string `json:"website"`
	UsesCount uint64 `json:"usesCount"`
}

// TopDomainsResponse structure
type TopDomainsResponseData struct {
	Website   string `json:"website"`
	UsesCount uint64 `json:"usesCount"`
}

// ShortlinkResponse structure
type ShortlinkResponse struct {
	Data   ShortlinkResponseData `json:"data"`
	Result string                `json:"result"`
}

// ShortlinkFullResponse structure
type ShortlinkFullResponse struct {
	Data   Shortlink `json:"data"`
	Result string    `json:"result"`
}

// ShortlinksResponse structure
type ShortlinksResponse struct {
	Data   []ShortlinkResponseData `json:"data"`
	Result string                  `json:"result"`
}

// TopDomainsResponse structure
type TopDomainsResponse struct {
	Data   []TopDomainsResponseData `json:"data"`
	Result string                   `json:"result"`
}

// DataMinutes : Information about uses per minute
type DataMinutes = map[int]int

// DataHours : Information about uses per hour
type DataHours = map[int]DataMinutes

// ShortlinksGraphResponseData : List of uses per day/hour/minute
type ShortlinksGraphResponseData = map[string]DataHours

// ShortlinksGraphResponse : Information about uses per day/hour/minute
type ShortlinksGraphResponse struct {
	Data   ShortlinksGraphResponseData `json:"data"`
	Result string                      `json:"result"`
}

// AddUseToGraph : Something
func AddUseToGraph(g *ShortlinksGraphResponseData, u ShortlinkUse) {
	formattedDate := u.UseTime.Format("2006-01-02")
	formattedHour := u.UseTime.Hour()
	formattedMinute := u.UseTime.Minute()

	if _, exists := (*g)[formattedDate]; !exists {
		(*g)[formattedDate] = make(DataHours)
		(*g)[formattedDate][formattedHour] = make(DataMinutes)
		(*g)[formattedDate][formattedHour][formattedMinute] = 1
	} else {
		if _, exists := (*g)[formattedDate][formattedHour]; !exists {
			(*g)[formattedDate][formattedHour] = make(DataMinutes)
			(*g)[formattedDate][formattedHour][formattedMinute] = 1
		} else {
			if _, exists := (*g)[formattedDate][formattedHour][formattedMinute]; !exists {
				(*g)[formattedDate][formattedHour][formattedMinute] = 1
			} else {
				(*g)[formattedDate][formattedHour][formattedMinute]++
			}
		}
	}
}
