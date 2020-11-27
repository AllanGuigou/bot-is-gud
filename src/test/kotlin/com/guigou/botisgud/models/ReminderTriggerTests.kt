package com.guigou.botisgud.models

import assertk.assertThat
import assertk.assertions.isEqualTo
import assertk.assertions.isGreaterThan
import org.junit.jupiter.api.*
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
            val sut = RelativeReminderTrigger(1, ChronoUnit.HOURS)

            val result = sut.timestamp(clock)

            assertThat(result).isEqualTo(Instant.parse("2020-08-08T11:15:30Z"))
        }
    }

    @Nested
    inner class `AbsoluteReminderTrigger` {
        @TestFactory
        fun `timestamp returns Instant of next UTC 5 PM`() =
            listOf(
                "2020-08-08T10:15:30Z" to "2020-08-08T17:00:00Z",
                "2020-08-08T17:15:30Z" to "2020-08-09T17:00:00Z",
                "2020-08-08T19:15:30Z" to "2020-08-09T17:00:00Z",
            ).map { (input, expected) ->
                DynamicTest.dynamicTest("when given $input") {
                    val clock = Clock.fixed(Instant.parse(input), ZoneId.of("UTC"))
                    val sut = AbsoluteReminderTrigger("0 17 * * *")

                    val result = sut.timestamp(clock)

                    assertThat(result).isEqualTo(Instant.parse(expected))
                }
            }

        @TestFactory
        fun `timestamp returns Instant of next America New York 5 PM `() =
            listOf(
                "2020-08-08T10:15:30Z" to "2020-08-08T21:00:00Z", // EDT
                "2020-08-08T21:15:30Z" to "2020-08-09T21:00:00Z", // EDT
                "2020-12-08T21:15:30Z" to "2020-12-08T22:00:00Z", // EST
                "2020-12-08T22:15:30Z" to "2020-12-09T22:00:00Z", // EST
            ).map { (input, expected) ->
                DynamicTest.dynamicTest("when given $input") {
                    val clock = Clock.fixed(Instant.parse(input), ZoneId.of("America/New_York"))
                    val sut = AbsoluteReminderTrigger("0 17 * * *")

                    val result = sut.timestamp(clock)

                    assertThat(result).isEqualTo(Instant.parse(expected))
                }
            }

        @TestFactory
        fun `constructor throws exception when given invalid input`() =
            listOf(
                "",
                "* * * *",
                "a 17 * * *",
            ).map { input ->
                DynamicTest.dynamicTest("when given $input") {
                    assertThrows<IllegalArgumentException> { AbsoluteReminderTrigger(input) }
                }
            }
    }
}
