package main

type Account struct {
	AccessToken string
	TokenType   string
	ExpiresIn   string
	Scope       string
}

type KordisResponseProfile struct {
	ResponseCode int    `json:"response_code"`
	Version      string `json:"version"`
	Result       struct {
		UID              int         `json:"uid"`
		StudentID        string      `json:"student_id"`
		Ine              string      `json:"ine"`
		Civility         string      `json:"civility"`
		Firstname        string      `json:"firstname"`
		Name             string      `json:"name"`
		MaidenName       interface{} `json:"maiden_name"`
		Birthday         int64       `json:"birthday"`
		Birthplace       string      `json:"birthplace"`
		BirthCountry     string      `json:"birth_country"`
		Address1         string      `json:"address1"`
		Address2         interface{} `json:"address2"`
		City             string      `json:"city"`
		Zipcode          string      `json:"zipcode"`
		Country          string      `json:"country"`
		Telephone        interface{} `json:"telephone"`
		Mobile           string      `json:"mobile"`
		Email            string      `json:"email"`
		Nationality      string      `json:"nationality"`
		PersonalMail     string      `json:"personal_mail"`
		Mailing          interface{} `json:"mailing"`
		EmergencyContact struct {
			EmergencyID int         `json:"emergency_id"`
			Type        interface{} `json:"type"`
			TypeDetails interface{} `json:"type_details"`
			Firstname   interface{} `json:"firstname"`
			Name        interface{} `json:"name"`
			Telephone   interface{} `json:"telephone"`
			Mobile      interface{} `json:"mobile"`
			WorkPhone   interface{} `json:"work_phone"`
		} `json:"emergency_contact"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Years struct {
				Href string `json:"href"`
			} `json:"years"`
			Agenda struct {
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"agenda"`
			Grades struct {
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"grades"`
			Classes struct {
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"classes"`
			Courses struct {
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"courses"`
			Teachers struct {
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"teachers"`
			News struct {
				Href string `json:"href"`
			} `json:"news"`
			Photo struct {
				Href string `json:"href"`
			} `json:"photo"`
		} `json:"_links"`
	} `json:"result"`
}

type KordisResultAgenda struct {
	ReservationID int `json:"reservation_id"`
	Rooms         []struct {
		Links     []interface{} `json:"links"`
		RoomID    int           `json:"room_id"`
		Name      string        `json:"name"`
		Floor     string        `json:"floor"`
		Campus    string        `json:"campus"`
		Color     string        `json:"color"`
		Latitude  string        `json:"latitude"`
		Longitude string        `json:"longitude"`
	} `json:"rooms"`
	Type       string      `json:"type"`
	Modality   string      `json:"modality"`
	Author     int         `json:"author"`
	CreateDate interface{} `json:"create_date"`
	StartDate  int64       `json:"start_date"`
	EndDate    int64       `json:"end_date"`
	State      string      `json:"state"`
	Comment    interface{} `json:"comment"`
	Classes    interface{} `json:"classes"`
	Name       string      `json:"name"`
	Discipline struct {
		Coef             interface{}   `json:"coef"`
		Ects             interface{}   `json:"ects"`
		Name             string        `json:"name"`
		Teacher          string        `json:"teacher"`
		Trimester        string        `json:"trimester"`
		Year             int           `json:"year"`
		Links            []interface{} `json:"links"`
		HasDocuments     interface{}   `json:"has_documents"`
		HasGrades        interface{}   `json:"has_grades"`
		NbStudents       int           `json:"nb_students"`
		RcID             int           `json:"rc_id"`
		SchoolID         int           `json:"school_id"`
		StudentGroupID   int           `json:"student_group_id"`
		StudentGroupName string        `json:"student_group_name"`
		SyllabusID       interface{}   `json:"syllabus_id"`
		TeacherID        int           `json:"teacher_id"`
		TrimesterID      int           `json:"trimester_id"`
	} `json:"discipline"`
	Teacher               string        `json:"teacher"`
	Promotion             string        `json:"promotion"`
	PrestationType        int           `json:"prestation_type"`
	IsElectronicSignature bool          `json:"is_electronic_signature"`
	Links                 []interface{} `json:"links"`
}
type KordisResponseAgenda struct {
	ResponseCode int                  `json:"response_code"`
	Version      string               `json:"version"`
	Result       []KordisResultAgenda `json:"result"`
	Links        []interface{}        `json:"links"`
}
