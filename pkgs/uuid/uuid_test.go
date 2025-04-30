package uuid

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewGenerator(t *testing.T) {
	Convey("Given a call to NewGenerator", t, func() {
		Convey("When creating a new Generator", func() {
			generator := NewGenerator()

			Convey("Then the generator should not be nil", func() {
				So(generator, ShouldNotBeNil)
			})

			Convey("Then the generator's New method should return a valid UUID string", func() {
				uuid := generator.New()

				// UUID v4 should be 36 characters (32 hex digits + 4 hyphens)
				So(len(uuid), ShouldEqual, 36)

				// Check that it matches the UUID format (8-4-4-4-12 hex digits)
				uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
				So(uuidPattern.MatchString(uuid), ShouldBeTrue)

				// Generate another UUID and ensure it's different (uniqueness check)
				uuid2 := generator.New()
				So(uuid2, ShouldNotEqual, uuid)
			})
		})
	})
}
