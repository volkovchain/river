# River


### Build

```shell
go build -o river
```

### Add user
```sql
INSERT INTO main.employers (name, addr, amount_salary) VALUES ('Firstname Lastname', '0xWalletAddress', 1000000);
```
### 1000000 amount salary? WTF!?

The example token we're using, USDC, uses 6 decimals which is standard practice for ERC-20 tokens. This means that in order to represent 1 token we have to do the calculation amount * 10^18. In this example we'll use 1 token so we'll need to calculate 1 * 10^6 which is 1000000

### Start pay a salary

```shell
./river
```

### Repay

```shell
./river repay
```