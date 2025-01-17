insert into collector_portfolios (totalamountshares,
                                  totalamountbonds,
                                  totalamountetf,
                                  totalamountcurrencies,
                                  totalamountfutures,
                                  expectedyield,
                                  accountid,
                                  totalamountoptions,
                                  totalamountsp,
                                  totalamountportfolio)
values (:totalamountshares,
        :totalamountbonds,
        :totalamountetf,
        :totalamountcurrencies,
        :totalamountfutures,
        :expectedyield,
        :accountid,
        :totalamountoptions,
        :totalamountsp,
        :totalamountportfolio)
on conflict (accountid) do nothing;