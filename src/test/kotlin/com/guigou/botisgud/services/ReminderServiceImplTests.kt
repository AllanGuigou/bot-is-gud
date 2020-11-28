package com.guigou.botisgud.services

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.models.RelativeReminderTrigger
import com.guigou.botisgud.models.ReminderDto
import io.ktor.http.Url
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.take
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Assertions.assertNotNull
import org.junit.jupiter.api.Disabled
import org.junit.jupiter.api.Test
import java.time.temporal.ChronoUnit

class ReminderServiceImplTests {

    @Test
    @Disabled
    fun `get returns reminders`() = runBlocking {
        // https://medium.com/@heyitsmohit/unit-testing-delays-errors-retries-with-kotlin-flows-77ce00d0c2f3
        val trigger = RelativeReminderTrigger(1, ChronoUnit.MILLIS)
        val reminder = ReminderDto(Snowflake("1"), "foo", Url("https://example.com"), trigger)
        val sut = ReminderServiceImpl()
        sut.add(reminder)

        val flow = sut.get()

        flow.take(5).collect { value -> assertNotNull(value) }
    }
}
