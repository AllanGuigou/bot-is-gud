package com.guigou.botisgud.services

import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import com.guigou.botisgud.models.timestamp
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.flow
import java.time.Instant

class ReminderServiceImpl : ReminderService {
    companion object {
        val logger = logger()
    }

    private val reminders = mutableListOf<Reminder>()

    override fun add(reminderDto: ReminderDto, trigger: ReminderTrigger) {
        val zoneId = "America/New_York"
        val timestamp = trigger.timestamp(zoneId)
        val reminder = Reminder(reminderDto.userId, reminderDto.message, reminderDto.link, timestamp)
        reminders.add(reminder)
        logger.info("add reminder for [${reminderDto.userId}] at [${timestamp}]")
    }

    override fun remove(reminder: Reminder) {
        reminders.remove(reminder)
        logger.info("remove reminder for [${reminder.userId}]")
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
