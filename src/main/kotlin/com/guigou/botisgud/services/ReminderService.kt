package com.guigou.botisgud.services

import com.guigou.botisgud.models.Reminder
import kotlinx.coroutines.flow.Flow

interface ReminderService {
    fun add(reminder: Reminder)
    fun remove(reminder: Reminder)
    suspend fun get(): Flow<Reminder>
}
