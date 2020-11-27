package com.guigou.botisgud.models

import com.cronutils.model.CronType
import com.cronutils.model.definition.CronDefinitionBuilder
import com.cronutils.model.time.ExecutionTime
import com.cronutils.parser.CronParser
import com.guigou.botisgud.extensions.java.util.unwrap
import java.time.Clock
import java.time.Instant
import java.time.ZoneId
import java.time.ZoneId.SHORT_IDS
import java.time.ZonedDateTime
import java.time.temporal.TemporalUnit

interface ReminderTrigger {
    fun timestamp(clock: Clock = Clock.systemUTC()): Instant
}

class RelativeReminderTrigger(
    private val value: Long,
    private val unit: TemporalUnit,
) : ReminderTrigger {
    override fun timestamp(clock: Clock): Instant {
        return Instant.now(clock)
            .plus(value, unit)
    }
}

class AbsoluteReminderTrigger(expression: String) : ReminderTrigger {
    private var time: ExecutionTime
    private val parser = CronParser(CronDefinitionBuilder.instanceDefinitionFor(CronType.UNIX))

    init {
        time = ExecutionTime.forCron(parser.parse(expression))
    }

    override fun timestamp(clock: Clock): Instant {
        val now = ZonedDateTime.now(clock)

        return time.nextExecution(now) // get next execution in the current zone
            .unwrap()!! // TODO: when would this be null? how should null be handled?
            .toInstant() // convert to Instant and UTC
    }
}

fun ReminderTrigger.timestamp(zoneId: String) = this.timestamp(Clock.system(ZoneId.of(zoneId)))
