package com.guigou.botisgud.services

import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import kotlinx.coroutines.flow.Flow

interface ReminderService {
    fun add(reminderDto: ReminderDto, trigger: ReminderTrigger)
    fun remove(reminder: Reminder)
    suspend fun get(): Flow<Reminder>
}
