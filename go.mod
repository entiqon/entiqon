module github.com/entiqon/entiqon

go 1.24

require (
    github.com/entiqon/common v1.0.0 // indirect until you bump
    github.com/entiqon/db v1.0.0     // indirect until you bump
)

replace github.com/entiqon/common => ./common
replace github.com/entiqon/db => ./db
