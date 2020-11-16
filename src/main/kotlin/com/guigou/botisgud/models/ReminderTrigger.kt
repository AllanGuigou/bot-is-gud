package com.guigou.botisgud.models

import com.cronutils.builder.CronBuilder
import com.cronutils.model.CronType
import com.cronutils.model.definition.CronDefinition
import com.cronutils.model.definition.CronDefinitionBuilder
import com.cronutils.model.field.expression.FieldExpressionFactory.on
import com.cronutils.model.time.ExecutionTime
import com.cronutils.parser.CronParser
import java.time.Clock
import java.time.Instant
import java.time.ZonedDateTime
import java.time.temporal.TemporalUnit
import java.util.*
import kotlin.reflect.typeOf

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

class AbsoluteReminderTrigger(private val clock: Clock = Clock.systemUTC()) : ReminderTrigger {
    @ExperimentalStdlibApi
    override fun timestamp(): Instant {
        val now = ZonedDateTime.now(clock)
        val parser = CronParser(CronDefinitionBuilder.instanceDefinitionFor(CronType.UNIX))
        val result = ExecutionTime.forCron(parser.parse("* 17 * * *")).nextExecution(now)

        return result.unwrap()!!.toInstant()
    }
}

fun <T> Optional<T>.unwrap(): T? = orElse(null)
