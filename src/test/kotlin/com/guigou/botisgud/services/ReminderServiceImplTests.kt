package com.guigou.botisgud.services

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.models.RelativeReminderTrigger
import com.guigou.botisgud.models.ReminderDto
import com.guigou.botisgud.services.reminder.ReminderServiceImpl
import com.guigou.botisgud.services.reminder.Reminders
import com.guigou.botisgud.services.user.Users
import io.ktor.http.*
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.take
import kotlinx.coroutines.runBlocking
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils.create
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.junit.jupiter.api.Assertions.assertNotNull
import org.junit.jupiter.api.Test
import java.time.temporal.ChronoUnit

class ReminderServiceImplTests {

    init {
        Database.connect("jdbc:sqlite:file:test?mode=memory&cache=shared", "org.sqlite.JDBC")
    }

    @Test
    fun `get returns reminders`() = runBlocking {
        // https://medium.com/@heyitsmohit/unit-testing-delays-errors-retries-with-kotlin-flows-77ce00d0c2f3
        val trigger = RelativeReminderTrigger(1, ChronoUnit.MILLIS)
        val reminder = ReminderDto(Snowflake("1"), "foo", Url("https://example.com"), trigger)
        val sut = ReminderServiceImpl()
        newSuspendedTransaction {
            create(Users)
            create(Reminders)

            sut.add(reminder)

            val flow = sut.get()

            flow.take(5).collect { value -> assertNotNull(value) }
        }
    }
}
