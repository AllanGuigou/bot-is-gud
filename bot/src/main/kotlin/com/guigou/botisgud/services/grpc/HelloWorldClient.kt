package com.guigou.botisgud.services.grpc

import io.grpc.ManagedChannel
import io.grpc.StatusException
import java.io.Closeable
import java.util.concurrent.TimeUnit
import kotlinx.coroutines.runBlocking
import com.example.*

/**
  Usage:
    val builder = ManagedChannelBuilder.forTarget("localhost:50051").usePlaintext()

    HelloWorldClient(
    // what was the point of an executor?
    builder.build()
    ).use {
    val user = "world"
    it.greet(user)
    }
 */
class HelloWorldClient(val channel: ManagedChannel) : Closeable {

    private val stub: GreeterGrpcKt.GreeterCoroutineStub = GreeterGrpcKt.GreeterCoroutineStub(channel)

    fun greet(s: String) = runBlocking {
        val request = helloRequest { name = s }
        try {
            val response = stub.sayHello(request)
            println("Greeter client received: ${response.message}")
        } catch (e: StatusException) {
            println("RPC failed: ${e.status}")
        }
    }

    override fun close() {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}