using System;
using System.Threading.Tasks;
using Grpc.Core;
using Presence;

class Program
{
	public static ManualResetEvent Shutdown = new ManualResetEvent(false);
	const int Port = 50051;

	public static void Main(string[] args)
	{
		Server server = new Server
		{
			Services = { PresenceService.BindService(new PresenceServiceImpl()) },
			Ports = { new ServerPort("localhost", Port, ServerCredentials.Insecure) }
		};

		server.Start();

		Console.WriteLine("Presence server listening on port " + Port);

		Shutdown.WaitOne();

		server.ShutdownAsync().Wait();
	}
}