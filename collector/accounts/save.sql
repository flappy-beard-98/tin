insert into collector_accounts(id,
                               type,
                               name,
                               status,
                               openeddate,
                               closeddate,
                               accesslevel)
values (:id,
        :type,
        :name,
        :status,
        :openeddate,
        :closeddate,
        :accesslevel)
on conflict (id) do nothing
;