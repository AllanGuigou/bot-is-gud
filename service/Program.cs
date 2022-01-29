﻿using System;
using System.Threading.Tasks;
using Grpc.Core;
using Helloworld;

class GreeterImpl : Greeter.GreeterBase
{
	public override Task<HelloReply> SayHello(HelloRequest request, ServerCallContext context)
	{
		Console.WriteLine("Hello " + request.Name);

		return Task.FromResult(new HelloReply { Message = "Hello " + request.Name });
	}
}

class Program
{
	const int Port = 50051;

	public static void Main(string[] args)
	{
		Server server = new Server
		{
			Services = { Greeter.BindService(new GreeterImpl()) },
			Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
		};

		server.Start();

		Console.WriteLine("Greeter server listening on port " + Port);
		Console.WriteLine("Press any key to stop the server...");
		Console.ReadKey();

		server.ShutdownAsync().Wait();
	}
}