-- +goose up

ALter table bookings alter column ID set not null;
ALter table bookings alter column user_id set not null;
ALter table bookings alter column pet_id set not null;

ALter table pet_status_log alter column ID set not null;
ALter table pet_status_log alter column id_bookings set not null;

ALter table notifications alter column ID set not null;
ALter table notifications alter column id_user set not null;

ALter table message alter column ID set not null;
ALter table message alter column bookings_id set not null;
ALter table message alter column receiver_id set not null;
ALter table message alter column sender_id set not null;




-- +goose down

ALter table bookings alter column ID DROP not null;
ALter table bookings alter column user_id DROP  not null;
ALter table bookings alter column pet_id DROP not null;

ALter table pet_status_log alter column ID DROP not null;
ALter table pet_status_log alter column id_bookings DROP  not null;

ALter table notifications alter column ID DROP  not null;
ALter table notifications alter column id_user DROP  not null;

ALter table message alter column ID DROP not null;
ALter table message alter column bookings_id DROP  not null;
ALter table message alter column receiver_id DROP  not null;
ALter table message alter column sender_id DROP not null;