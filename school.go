package gofakeit

import "math/rand/v2"

// School will generate a random School type
func School() string { return school(GlobalFaker.Rand) }

func (f *Faker) School() string { return school(f.Rand) }

func school(r *rand.Rand) string {
	return getRandValue(
		r, []string{"school", "name"}) + " " +
		getRandValue(r, []string{"school", "isPrivate"}) + " " +
		getRandValue(r, []string{"school", "type"})
}

func addSchoolLookup() {
	AddFuncLookup("school", Info{
		Display:     "School",
		Category:    "school",
		Description: "An institution for formal education and learning",
		Example:     `Harborview State Academy`,
		Output:      "string",
		Generate: func(r *rand.Rand, m *MapParams, info *Info) (any, error) {
			return school(r), nil
		},
	})
}
