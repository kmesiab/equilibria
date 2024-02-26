select t.name,
       t.bill_rate_in_credits,
       m.reference_id,
       s.name                                     as status,
       CONCAT(u1.firstname, ' ', u1.phone_number) as 'from',
       CONCAT(u2.firstname, ' ', u2.phone_number) as 'to',
       m.body,
       m.sent_at,
       m.received_at

from messages m
         inner join users u1 on m.from_user_id = u1.id
         inner join users u2 on m.to_user_id = u2.id
         inner join message_statuses s on s.id = m.message_status_id
         inner join message_types t on m.message_type_id = t.id
         inner join conversations c on m.conversation_id = c.id
where from_user_id = 4
   or to_user_id = 4;

select sum(bill_rate_in_credits) as credits_used

from messages m
         inner join users u1 on m.from_user_id = u1.id
         inner join users u2 on m.to_user_id = u2.id
         inner join message_statuses s on s.id = m.message_status_id
         inner join message_types t on m.message_type_id = t.id
         inner join conversations c on m.conversation_id = c.id
where to_user_id = 4
  and s.name = 'delivered';
