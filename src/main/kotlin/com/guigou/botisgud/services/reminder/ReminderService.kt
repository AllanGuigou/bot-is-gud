package com.guigou.botisgud.services.reminder

import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import kotlinx.coroutines.flow.Flow

interface ReminderService {
    suspend fun add(dto: ReminderDto): Reminder
    suspend fun remove(id: Int)
    suspend fun get(): Flow<Reminder>
}
