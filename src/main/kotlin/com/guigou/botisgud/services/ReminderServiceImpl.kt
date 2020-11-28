package com.guigou.botisgud.services

import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.Reminder
import com.guigou.botisgud.models.ReminderDto
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.flow
import org.litote.kmongo.Id
import org.litote.kmongo.coroutine.CoroutineCollection
import org.litote.kmongo.coroutine.coroutine
import org.litote.kmongo.eq
import org.litote.kmongo.lte
import org.litote.kmongo.reactivestreams.KMongo
import java.time.Instant

class ReminderServiceImpl : ReminderService {
    companion object {
        val logger = logger()
    }

    private val userService: UserService
    private val reminders: CoroutineCollection<Reminder>

    init {
        val client = KMongo.createClient().coroutine
        val database = client.getDatabase("bot-is-gud")
        userService = UserServiceImpl(database)
        reminders = database.getCollection("reminders")
    }

    override suspend fun add(reminderDto: ReminderDto): Reminder {
        val zone = userService.get(reminderDto.userId)?.zone ?: "America/New_York"

        reminderDto.toReminder(zone).let {
            reminders.insertOne(it)
            logger.info("add reminder [${it.key}] for [${it.userId}] at [${it.timestamp}]")
            return it
        }
    }

    override suspend fun remove(key: Id<Reminder>) {
        reminders.deleteOne(Reminder::key eq key)
        logger.info("remove reminder [${key}]")
    }

    override suspend fun get() = flow {
        // TODO: investigate if infinite flows are kosher
        while (true) {
            reminders.find(Reminder::timestamp lte Instant.now())
                .toFlow().collect {
                    emit(it)
                }

            delay(1000)
        }
    }
}
