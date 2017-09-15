//coupon_test provides unit tests for business domain model of coupon
package coupon_test

import (
	"fmt"
	"os"
	"sstest/model/coupon"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

var newCoupon, activeCoupon, inactiveCoupon, suspendedCoupon, noStockCoupon, earlyCoupon, expiredCoupon *coupon.Coupon
var currentTime = time.Now()
var beginningOfToday = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
var beginningOfTomorrow = beginningOfToday.AddDate(0, 0, 1)
var beginningOf90DaysFromToday = beginningOfToday.AddDate(0, 0, 90)
var beginningOfYesterday = beginningOfToday.AddDate(0, -1, 0)
var beginningOfLastWeek = beginningOfToday.AddDate(0, -7, 0)

func TestMain(m *testing.M) {
	//test setup
	newCoupon = coupon.New("newCoupon")

	activeCoupon = coupon.New("activeCoupon")
	activeCoupon.SetStatus(coupon.StatusActive)
	activeCoupon.SetStock(100)
	activeCoupon.SetKind(coupon.KindValue)
	activeCoupon.SetValue(decimal.New(10000, 0))

	inactiveCoupon = coupon.New("inactiveCoupon")
	inactiveCoupon.SetStatus(coupon.StatusInactive)
	inactiveCoupon.SetStock(50)
	inactiveCoupon.SetKind(coupon.KindPercentage)
	inactiveCoupon.SetValue(decimal.New(50, 0))

	suspendedCoupon = coupon.New("suspendedCoupon")
	suspendedCoupon.SetStatus(coupon.StatusSuspended)
	suspendedCoupon.SetStock(10)
	suspendedCoupon.SetKind(coupon.KindValue)
	suspendedCoupon.SetValue(decimal.New(25, 0))

	noStockCoupon = coupon.New("noStockCoupon")
	noStockCoupon.SetStatus(coupon.StatusActive)
	noStockCoupon.SetStock(0)
	noStockCoupon.SetKind(coupon.KindPercentage)
	noStockCoupon.SetValue(decimal.New(30, 0))

	earlyCoupon = coupon.New("earlyCoupon")
	earlyCoupon.SetStatus(coupon.StatusActive)
	earlyCoupon.SetStock(20)
	earlyCoupon.SetStartDate(beginningOfTomorrow)
	earlyCoupon.SetEndDate(beginningOf90DaysFromToday)
	earlyCoupon.SetKind(coupon.KindValue)
	earlyCoupon.SetValue(decimal.New(10000, 0))

	expiredCoupon = coupon.New("expiredCoupon")
	expiredCoupon.SetStatus(coupon.StatusActive)
	expiredCoupon.SetStock(20)
	expiredCoupon.SetStartDate(beginningOfLastWeek)
	expiredCoupon.SetEndDate(beginningOfYesterday)
	expiredCoupon.SetKind(coupon.KindPercentage)
	expiredCoupon.SetValue(decimal.New(80, 0))

	//run tests
	exitCode := m.Run()
	//test teardown
	os.Exit(exitCode)
}

func TestNewlyCreatedCoupon(t *testing.T) {
	t.Run("Status Must Be Inactive", func(t *testing.T) {
		if coupon.StatusInactive != newCoupon.Status() {
			t.Errorf("expected %v but got %v", coupon.StatusInactive, newCoupon.Status())
		}
	})
	t.Run("Stock Must Be Zero", func(t *testing.T) {
		if 0 != newCoupon.Stock() {
			t.Errorf("expected %v but got %v", 0, newCoupon.Stock())
		}
	})
	t.Run("Kind Must Be Percentage", func(t *testing.T) {
		if coupon.KindPercentage != newCoupon.Kind() {
			t.Errorf("expected %v but got %v", coupon.KindPercentage, newCoupon.Kind())
		}
	})
	t.Run("Value Must Be 10", func(t *testing.T) {
		if false == decimal.New(10, 0).Equal(newCoupon.Value()) {
			t.Errorf("expected %v but got %v", decimal.New(10, 0).String(), newCoupon.Value().String())
		}
	})
	t.Run("Start Date Must Be Beginning Of Today", func(t *testing.T) {
		if false == beginningOfToday.Equal(newCoupon.StartDate()) {
			t.Errorf("coupon start date %v does not match %v", newCoupon.StartDate(), beginningOfToday)
		}
	})
	t.Run("End Date Must Be Beginning Of 90 Days From Today", func(t *testing.T) {
		if false == beginningOf90DaysFromToday.Equal(newCoupon.EndDate()) {
			t.Errorf("coupon end date %v does not match %v", newCoupon.EndDate(), beginningOf90DaysFromToday)
		}
	})
}

func TestSettingUnknownStatusValue(t *testing.T) {
	_, err := newCoupon.SetStatus("UnknownStatus")
	if err == nil {
		t.Error("expected error but got none\n")
	}
}
func TestSettingNegativeValue(t *testing.T) {
	_, err := newCoupon.SetStatus("UnknownStatus")
	if err == nil {
		t.Error("expected error but got none\n")
	}
}

func TestSettingValue(t *testing.T) {
	kindValueCoupon1 := coupon.New("kindValue1")
	kindValueCoupon1.SetKind(coupon.KindValue)
	kindValueCoupon1.SetValue(decimal.New(10000, 0))

	kindPercentageCoupon1 := coupon.New("kindPercentage1")
	kindPercentageCoupon1.SetKind(coupon.KindPercentage)
	kindPercentageCoupon1.SetValue(decimal.New(20, 0))

	kindValueCoupon2 := coupon.New("kindValue2")
	kindValueCoupon2.SetKind(coupon.KindValue)
	kindValueCoupon2.SetValue(decimal.New(10000, 0))

	kindPercentageCoupon2 := coupon.New("kindPercentage2")
	kindPercentageCoupon2.SetKind(coupon.KindPercentage)
	kindPercentageCoupon2.SetValue(decimal.New(20, 0))

	kindValueCoupon3 := coupon.New("kindValue3")
	kindValueCoupon3.SetKind(coupon.KindValue)
	kindValueCoupon3.SetValue(decimal.New(5000, 0))

	_, errSetValueZeroOnKindValue := kindValueCoupon1.SetValue(decimal.New(0, 0))
	_, errSetValueZeroOnKindPercentage := kindPercentageCoupon1.SetValue(decimal.New(0, 0))
	_, errSetValueAbove100OnKindValue := kindValueCoupon2.SetValue(decimal.New(3000, 0))
	_, errSetValueAbove100OnKindPercentage := kindPercentageCoupon1.SetValue(decimal.New(3000, 0))
	_, errSetValue50OnKindValue := kindValueCoupon3.SetValue(decimal.New(50, 0))
	_, errSetValue50OnKindPercentage := kindPercentageCoupon2.SetValue(decimal.New(50, 0))

	var settingValueTests = []struct {
		testCase         string
		expectedValue    decimal.Decimal
		actualValue      decimal.Decimal
		expectedErrIsNil bool
		actualErrIsNil   bool
	}{
		{"Setting Value Zero On Kind Value", decimal.New(10000, 0), kindValueCoupon1.Value(), false, errSetValueZeroOnKindValue == nil},
		{"Setting Value Zero On Kind Percentage", decimal.New(20, 0), kindPercentageCoupon1.Value(), false, errSetValueZeroOnKindPercentage == nil},
		{"Setting Value Above 100 On Kind Value", decimal.New(3000, 0), kindValueCoupon2.Value(), true, errSetValueAbove100OnKindValue == nil},
		{"Setting Value Above 100 On Kind Percentage", decimal.New(20, 0), kindPercentageCoupon1.Value(), false, errSetValueAbove100OnKindPercentage == nil},
		{"Setting Value 50 On Kind Value", decimal.New(50, 0), kindValueCoupon3.Value(), true, errSetValue50OnKindValue == nil},
		{"Setting Value 50 On Kind Percentage", decimal.New(50, 0), kindPercentageCoupon2.Value(), true, errSetValue50OnKindPercentage == nil},
	}

	for _, test := range settingValueTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if false == test.expectedValue.Equal(test.actualValue) {
				t.Errorf("want %v for value, got %v", test.expectedValue.String(), test.actualValue.String())
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}
}

func TestSettingKind(t *testing.T) {
	value50CouponKindValue := coupon.New("value50")
	value50CouponKindValue.SetKind(coupon.KindValue)
	value50CouponKindValue.SetValue(decimal.New(50, 0))

	value50CouponKindPercentage := coupon.New("value50")
	value50CouponKindPercentage.SetKind(coupon.KindPercentage)
	value50CouponKindPercentage.SetValue(decimal.New(50, 0))

	value2000CouponKindValue := coupon.New("value2000")
	value2000CouponKindValue.SetKind(coupon.KindValue)
	value2000CouponKindValue.SetValue(decimal.New(2000, 0))

	_, errSetKindPercentageOn50Value := value50CouponKindValue.SetKind(coupon.KindPercentage)
	_, errSetKindValueOn50Value := value50CouponKindPercentage.SetKind(coupon.KindValue)
	_, errSetKindPercentageOn2000Value := value2000CouponKindValue.SetKind(coupon.KindPercentage)

	var settingKindTests = []struct {
		testCase         string
		expectedValue    string
		actualValue      string
		expectedErrIsNil bool
		actualErrIsNil   bool
	}{
		{"Setting Kind Percentage On Value 50", coupon.KindPercentage, value50CouponKindValue.Kind(), true, errSetKindPercentageOn50Value == nil},
		{"Setting Kind Value On Value 50", coupon.KindValue, value50CouponKindPercentage.Kind(), true, errSetKindValueOn50Value == nil},
		{"Setting Kind Percentage On Value 2000", coupon.KindValue, value2000CouponKindValue.Kind(), false, errSetKindPercentageOn2000Value == nil},
	}

	for _, test := range settingKindTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for kind, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}
}

func TestSettingWrongStartAndEndDate(t *testing.T) {
	datedCoupon := coupon.New("datedCoupon")
	t.Run("Setting End Date Earlier Than Start Date Must Return Error", func(t *testing.T) {
		if _, err := datedCoupon.SetEndDate(beginningOfYesterday); err == nil {
			t.Error("expected error but got none\n")
		}
	})
	datedCoupon.SetEndDate(beginningOfTomorrow)
	t.Run("Setting Start Date Later Than End Date Must Return Error", func(t *testing.T) {
		if _, err := datedCoupon.SetStartDate(beginningOf90DaysFromToday); err == nil {
			t.Error("expected error but got none\n")
		}
	})
}

func TestGetDiscountAmount(t *testing.T) {
	kindValueCoupon := coupon.New("kindValue")
	kindValueCoupon.SetKind(coupon.KindValue)
	kindValueCoupon.SetValue(decimal.New(5000, 0))

	kindPercentageCoupon := coupon.New("kindPercentage")
	kindPercentageCoupon.SetKind(coupon.KindPercentage)
	kindPercentageCoupon.SetValue(decimal.New(25, 0))

	amount := decimal.New(50000, 0)

	t.Run("Getting Discount Amount From Coupon Kind Value", func(t *testing.T) {
		if false == decimal.New(5000, 0).Equal(kindValueCoupon.GetDiscountAmount(amount)) {
			t.Errorf("want %v got %v", decimal.New(5000, 0), kindValueCoupon.GetDiscountAmount(amount))
		}
	})

	t.Run("Getting Discount Amount From Coupon Kind Percentage", func(t *testing.T) {
		if false == decimal.New(12500, 0).Equal(kindPercentageCoupon.GetDiscountAmount(amount)) {
			t.Errorf("want %v got %v", decimal.New(12500, 0), kindPercentageCoupon.GetDiscountAmount(amount))
		}
	})
}

func TestCanBeApplied(t *testing.T) {
	validCouponApplication, _ := activeCoupon.CanBeApplied()
	invalidCouponApplication, _ := inactiveCoupon.CanBeApplied()
	suspendedCouponApplication, _ := suspendedCoupon.CanBeApplied()
	earlyCouponApplication, _ := earlyCoupon.CanBeApplied()
	expiredCouponApplication, _ := expiredCoupon.CanBeApplied()
	noStockCouponApplication, _ := noStockCoupon.CanBeApplied()

	var canBeOrderedTests = []struct {
		testCase string
		expected bool
		actual   bool
	}{
		{"Valid Coupon Application", true, validCouponApplication},
		{"Invalid Coupon Application", false, invalidCouponApplication},
		{"Suspended Coupon Application", false, suspendedCouponApplication},
		{"Early Coupon Application", false, earlyCouponApplication},
		{"Expired Coupon Application", false, expiredCouponApplication},
		{"Out of Stock Coupon Application", false, noStockCouponApplication},
	}

	for _, test := range canBeOrderedTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expected != test.actual {
				t.Errorf("want %v got %v", test.expected, test.actual)
			}
		})
	}
}
