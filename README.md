# Montui 

The cli budgeting money TUI.

# Setup

Requires a .env file with
```
SECRET_ID=
SECRET_KEY=
```

filled in.
These can be obtained from https://gocardless.com/bank-account-data/
# Example
![Example image of Montui displaying categorised transactions](<example.png>)

# Walkthrough 

Press tab to change between pages
Accounts -> Transactions -> Categories

On the "Accounts page" press `n`, this will walk you through connecting a bank account

On the transactions page press `r`, this will perform a network fetch and grab transactions

Prefetched transactions will load instantly on new launches, fetching is always manual.

Scroll the transactions page using arrows, press enter to "categorise" a transaction. Use arrows again to select 
a category. Press enter to lock in your selection.

Pressing tab again takes you to the categories page.
Scroll the category page with arrows, press enter to open the "edit" page

Use the arrow keys to move between the name or color field to edit.
Press enter again to commit your changes.

On the rules page, press n to create new rules, enter to edit an existing rule, and d to delete a rule
Pressing r will rerun the categorisation system.

## Notes

The generated spec is wrong
transaction response is 
{
    "transactions": {
        "booked": []
    }
}

not 
{
    "booked": []
}
currency exchange is not an array

