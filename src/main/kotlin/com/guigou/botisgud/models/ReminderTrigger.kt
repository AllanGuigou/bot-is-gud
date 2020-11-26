package com.guigou.botisgud.models

import com.cronutils.model.CronType
import com.cronutils.model.definition.CronDefinitionBuilder
import com.cronutils.model.time.ExecutionTime
import com.cronutils.parser.CronParser
import com.guigou.botisgud.extensions.java.util.unwrap
import java.time.Clock
import java.time.Instant
import java.time.ZonedDateTime
import java.time.temporal.TemporalUnit

interface ReminderTrigger {
    fun timestamp(): Instant
}

class RelativeReminderTrigger(
    private val value: Long,
    private val unit: TemporalUnit,
    private val clock: Clock = Clock.systemUTC()
) : ReminderTrigger {
    override fun timestamp(): Instant {
        return Instant.now(clock)
            .plus(value, unit)
    }
}

class AbsoluteReminderTrigger(expression: String, private var clock: Clock = Clock.systemUTC()) : ReminderTrigger {
    private var time: ExecutionTime
    private val parser = CronParser(CronDefinitionBuilder.instanceDefinitionFor(CronType.UNIX))

    init {
        time = ExecutionTime.forCron(parser.parse(expression))
    }

    override fun timestamp(): Instant {
        val now = ZonedDateTime.now(clock)
        val result = time.nextExecution(now)

        // TODO: how should null be handled?
        return result.unwrap()!!.toInstant()
    }
}
