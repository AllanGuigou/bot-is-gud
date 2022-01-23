package com.guigou.botisgud.services.reminder

import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.models.timestamp
import com.guigou.botisgud.services.user.UserService
import com.guigou.botisgud.services.user.UserServiceImpl
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.flow
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import java.time.Instant
import java.time.ZoneOffset

class ReminderServiceImpl() : ReminderService {
    companion object {
        val logger = logger()
    }

    private val userService: UserService

    init {
        userService = UserServiceImpl()
    }

    override suspend fun add(dto: ReminderDto): Reminder {
        val zone = userService.get(dto.userId)?.zone ?: "America/New_York"
        val timestamp = dto.trigger.timestamp(zone).atOffset(ZoneOffset.UTC).toLocalDateTime()
        return newSuspendedTransaction {
            val reminder = ReminderEntity.new {
                userId = dto.userId.value
                message = dto.message
                link = dto.link.toString()
                // https://stackoverflow.com/questions/52264768/how-to-convert-from-instant-to-localdate
                this.timestamp = timestamp
            }
            logger.info("add reminder [${reminder.id}] for [${reminder.userId}] at [${reminder.timestamp}]")
            return@newSuspendedTransaction reminder.toReminder();
        }
    }

    override suspend fun remove(id: Int) {
        newSuspendedTransaction {
            ReminderEntity.findById(id)?.delete()
            logger.info("remove reminder [${id}]")
        }
    }

    override suspend fun get() = flow<Reminder> {
        // TODO: investigate if infinite flows are kosher
        // TODO: what are alternative ways to handle get/delete to avoid a duplicate reminder race condition
        while (true) {
            newSuspendedTransaction {
                return@newSuspendedTransaction ReminderEntity.find {
                    Reminders.timestamp lessEq Instant.now().atOffset(ZoneOffset.UTC).toLocalDateTime()
                }.toList()
            }.forEach {
                emit(it.toReminder())
            }

            delay(1000)
        }
    }
}
