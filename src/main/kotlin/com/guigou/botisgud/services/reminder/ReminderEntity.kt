package com.guigou.botisgud.services.reminder

import com.guigou.botisgud.models.Reminder
import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import java.time.ZoneOffset

class ReminderEntity(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<ReminderEntity>(Reminders)

    var userId by Reminders.userId
    var message by Reminders.message
    var link by Reminders.link
    var timestamp by Reminders.timestamp
}

fun ReminderEntity.toReminder() = Reminder(
    id.value, userId, message, link, timestamp.toInstant(ZoneOffset.UTC)
)
