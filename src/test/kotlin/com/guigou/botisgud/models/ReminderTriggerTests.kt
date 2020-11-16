package com.guigou.botisgud.models

import assertk.assertThat
import assertk.assertions.isEqualTo
import assertk.assertions.isGreaterThan
import org.junit.jupiter.api.DynamicTest
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestFactory
import java.time.Clock
import java.time.Instant
import java.time.ZoneId
import java.time.temporal.ChronoUnit

class ReminderTriggerTests {

    @Nested
    inner class `RelativeReminderTrigger` {
        @Test
        fun `timestamp gives a time in the future`() {
            // TODO: this is not a very robust test
            // the goal here is to validate RelativeReminderTrigger with the default clock
            val before = Instant.now()
            val sut = RelativeReminderTrigger(1, ChronoUnit.MILLIS)

            val result = sut.timestamp()

            assertThat(result).isGreaterThan(before)
        }

        @Test
        fun `timestamp gives a time one hour in the future`() {
            val clock = Clock.fixed(Instant.parse("2020-08-08T10:15:30Z"), ZoneId.of("UTC"))

            val sut = RelativeReminderTrigger(1, ChronoUnit.HOURS, clock)

            val result = sut.timestamp()

            assertThat(result).isEqualTo(Instant.parse("2020-08-08T11:15:30Z"))
        }
    }

    // TODO: 2020-08-08T17:15:30Z
    @Nested
    inner class `AbsoluteReminderTrigger` {
        @ExperimentalStdlibApi
        @TestFactory
        fun `timestamp returns an instant for the next 5 PM UTC`() =
            listOf(
                "2020-08-08T10:15:30Z" to "2020-08-08T17:00:00Z",
                "2020-08-08T19:15:30Z" to "2020-08-09T17:00:00Z"
            ).map { (input, expected) ->
                DynamicTest.dynamicTest("when it is $input") {
                    val clock = Clock.fixed(Instant.parse(input), ZoneId.of("UTC"))

                    val sut = AbsoluteReminderTrigger(clock)

                    val result = sut.timestamp()

                    assertThat(result).isEqualTo(Instant.parse(expected))
                }
            }
    }
}
