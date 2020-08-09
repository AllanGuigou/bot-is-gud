package com.guigou.botisgud.models

import assertk.assertThat
import assertk.assertions.isEqualTo
import assertk.assertions.isGreaterThan
import org.junit.Test
import java.time.Clock
import java.time.Instant
import java.time.ZoneId
import java.time.temporal.ChronoUnit

class RelativeReminderTriggerTests {

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
