select figi,
       uid
from collector_shares
where currency = :currency
;
