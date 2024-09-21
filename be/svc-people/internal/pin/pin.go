package pin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/ddd"
)

var egn_WEIGHTS = []int{2, 4, 8, 5, 10, 9, 7, 3, 6}

type egn struct {
	err         error
	origin      string
	DateOfBirth time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	Egn         string    `json:"egn"`
	Ok          bool      `json:"ok"`
	month       int
	day         int
	year        int
}

func (e *egn) Err() error {
	return e.err
}

func Parse(e string) (ddd.PinValidation, error) {
	eg := &egn{origin: strings.TrimSpace(e)}
	eg.Ok = eg.validate()
	if eg.Ok {
		return ddd.PinValidation{
			Dob: ddd.Dob{
				Year:  eg.year + 1900,
				Month: eg.month,
				Day:   eg.day,
			},
			Gender: eg.Gender,
		}, nil
	}
	return ddd.PinValidation{}, eg.err
}

func (o *egn) validate() bool {
	var err error
	if len(o.origin) != 10 {
		o.err = fmt.Errorf("string len not 10")
		return false
	}
	// if (len(o.or) != 10):
	// 	return False
	year, err := strconv.Atoi(o.origin[0:2])
	if err != nil {
		o.err = fmt.Errorf("cant convert year: %v", err)
		return false
	}
	mon, err := strconv.Atoi(o.origin[2 : 2+2])
	if err != nil {
		o.err = fmt.Errorf("cant convert month: %v", err)
		return false
	}
	day, err := strconv.Atoi(o.origin[4 : 4+2])
	if err != nil {
		o.err = fmt.Errorf("cant convert day: %v", err)
		return false
	}

	dob := time.Time{}

	o.year = year
	o.month = mon
	o.day = day
	if o.month > 40 {
		dob, err = checkDateFormat(year+2000, mon-40, day)
	} else if o.month > 20 {
		dob, err = checkDateFormat(year+1800, mon-20, day)
	} else {
		dob, err = checkDateFormat(year+1900, mon, day)

	}
	if err != nil {
		o.err = err
		return false
	}
	o.DateOfBirth = dob
	sexBit, err := strconv.Atoi(o.origin[8 : 8+1])
	if err == nil {
		// log.Printf("asdasdasDASDASDAsdasd: %v", err)
		o.Gender = "female"
		if (sexBit % 2) == 0 {
			o.Gender = "male"
		}
	}
	checksum, err := strconv.Atoi(o.origin[9 : 9+1])
	if err != nil {
		o.err = fmt.Errorf("incorrect checksum value")
		return false
	}
	egnsum := 0
	for i, w := range egn_WEIGHTS {
		// w := egn_WEIGHTS[i]
		ii, err := strconv.Atoi(o.origin[i : i+1])
		if err != nil {
			o.err = fmt.Errorf("incorrect checksum value conver")
			return false
		}
		mm := ii * w
		egnsum = egnsum + mm
		// log.Printf("%v %v, %v %v sum:%v", i, w, ii, mm, egnsum)
	}
	valid_checksum := egnsum % 11
	// log.Printf("crc %v %v", checksum, valid_checksum)
	if valid_checksum == 10 {
		valid_checksum = 0
	}
	if valid_checksum == checksum {
		o.Egn = o.origin
		return true
	}
	// for i in range(0, 9):
	// egnsum = egnsum + ( int(egn[i: i + 1]) * EGN_WEIGHTS[i])
	// valid_checksum = egnsum % 11
	// if valid_checksum == 10:
	// 	valid_checksum = 0;
	// if checksum == valid_checksum:
	// 	return True
	o.err = fmt.Errorf("incorrect calc checksum")

	return false
}

func checkDateFormat(year int, month int, day int) (time.Time, error) {
	t, err := time.Parse("2/1/2006", fmt.Sprintf("%d/%d/%d", day, month, year))
	return t, err
}
