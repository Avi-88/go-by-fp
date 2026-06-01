# Day 6 Notes

## What is the difference between errors.Is and errors.As?

errors.Is is used to determine the exact error and is a direct comparison whereas the errors.As is more of a type comparison throughout the entire error tree ( unwraps automatically to find if there is a match )

## What does %w do that %v does not?

%w preserves the error tree and has a Unwrap method which can  be used to identify the error further while %v merges the error text into a flat string so to identify error we can only use fragile string matching

## When would you use a custom error type vs errors.New vs fmt.Errorf?

The errors.New has the standard go error structure so if I require a custom format of error I would use a custom type of error. fmt.Errorf is only used for printing errors , any parent caller wont be able to determine if the call failed or succeeded if it is used ( not part of the error tree )

## What does it mean for an error to be "wrapped"?

When an error is returned to a parent caller who has its own error structure the final returned error wraps over the childs thrown error

## What surprised you?

(your answer here)
