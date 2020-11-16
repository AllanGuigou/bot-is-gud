package com.guigou.botisgud.services

import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.flow
import java.time.Instant

class ReminderServiceImpl : ReminderService {
    private val reminders = mutableListOf<Reminder>()

    override fun add(reminderDto: ReminderDto, trigger: ReminderTrigger) {
        val reminder = Reminder(reminderDto.userId, reminderDto.message, reminderDto.link, trigger.timestamp())
        reminders.add(reminder)
    }

    override fun remove(reminder: Reminder) {
        reminders.remove(reminder)
    }

    override suspend fun get() = flow {
        // TODO: investigate if infinite flows are kosher
        while (true) {
            for (reminder in reminders.filter { it.timestamp <= Instant.now() }) {
                emit(reminder)
            }

            delay(1000)
        }
    }
}
