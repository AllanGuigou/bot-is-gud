package com.guigou.botisgud.services

import com.guigou.botisgud.models.Reminder
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.flow
import java.util.*

class ReminderServiceImpl : ReminderService {
    private val reminders = mutableListOf<Reminder>()

    override fun add(reminder: Reminder) {
        reminders.add(reminder)
    }

    override fun remove(reminder: Reminder) {
        reminders.remove(reminder)
    }

    override suspend fun get() = flow {
        // TODO: investigate if infinite flows are kosher
        while (true) {
            for (reminder in reminders.filter { it.timestamp <= Date() }) {
                emit(reminder)
            }

            delay(1000)
        }
    }

}
