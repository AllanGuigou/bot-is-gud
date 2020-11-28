package com.guigou.botisgud.services

import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.ReminderTrigger
import kotlinx.coroutines.flow.Flow
import org.litote.kmongo.Id

interface ReminderService {
    suspend fun add(reminderDto: ReminderDto): Reminder
    suspend fun remove(key: Id<Reminder>)
    suspend fun get(): Flow<Reminder>
}
