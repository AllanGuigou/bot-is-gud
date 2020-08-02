package com.guigou.botisgud.services

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.models.Reminder
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.take
import kotlinx.coroutines.runBlocking
import org.junit.Test
import java.util.*
import kotlin.test.assertNotNull

class ReminderServiceImplTests {

    @Test
    fun `test`() = runBlocking {
        // https://medium.com/@heyitsmohit/unit-testing-delays-errors-retries-with-kotlin-flows-77ce00d0c2f3
        val sut = ReminderServiceImpl()
        sut.add(Reminder(Snowflake("foo"), "bar", "baz", Date()))

        val flow = sut.get()

        flow.take(5).collect { value -> assertNotNull(value) }
    }
}
