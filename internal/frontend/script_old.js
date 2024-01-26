let graphData = {"name":"","deps":[{"name":"Processor","deps":[{"name":"measured.NewPreProcessor","deps":[{"name":"preprocessor.NewPreprocessor","deps":[{"name":"external.New","deps":[{"name":"external.ExtractAcIDForPX"},{"name":"external.FromPath"}]},{"name":"measured.NewProcessor","deps":[{"name":"cacheProcessor.New","deps":[{"name":"measured.NewSchemaRepo","deps":[{"name":"schema.NewRepo","deps":[{"name":"schemaClient.NewSchemaServiceClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]}]}]},{"name":"measured.NewTimeProcessor","deps":[{"name":"exampleevent.NewProcessor","deps":[{"name":"message.NewProcessor","deps":[{"name":"trigger.NewSelector","deps":[{"name":"measured.NewRecipeRepo","deps":[{"name":"cache.New","deps":[{"name":"recipe.NewRecipeClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]},{"name":"cfg.Meta.Application","value":"px"},{"name":"cfg.GetCacheRecipesTime"},{"name":"cfg.GetCacheRecipesCleanupTime"}]}]},{"name":"cfg.Meta.ScenarioType","value":"px"}]},{"name":"measured.NewSingleProcessor","deps":[{"name":"single.NewProcessor","deps":[{"name":"readerfactory.New","deps":[{"name":"","deps":[{"name":"\"MAP\"","deps":[{"name":"dict.NewReaderFactory"}]},{"name":"\"EVENT\"","deps":[{"name":"eventReader.NewReaderFactory"}]},{"name":"\"SEGMENT\"","deps":[{"name":"segmentReader.NewReaderFactory","deps":[{"name":"segment.NewSegmentServiceClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]},{"name":"dp.NewDigitalProfileServiceClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]}]}]},{"name":"\"PROFILE\"","deps":[{"name":"profile.NewReaderFactory","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]}]}]}]},{"name":"delaythrottle.NewRepository","deps":[{"name":"delayedData.NewRepository","deps":[{"name":"collection.NewProvider","deps":[{"name":"mongoPkg.NewClient","deps":[{"name":"strings.ReplaceAll","deps":[{"name":"cfg.Delay.Mongo.ConnectionString","value":"mongodb://srv_cdp:{{PASSWORD}}@mongo-45-1.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-2.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-3.cdp.dc-2.dev.dcwp.pl:27071/?replicaSet=cdp\u0026readPreference=primary\u0026serverSelectionTimeoutMS=5000\u0026connectTimeoutMS=10000\u0026authSource=admin\u0026authMechanism=SCRAM-SHA-256"},{"name":"cfg.Delay.Mongo.Password"}]},{"name":"cfg.GetMongoTimeout"}]},{"name":"cfg.Delay.Mongo.DB","value":"delay"},{"name":"cfg.Delay.Mongo.DelayedDataCollection","value":"data"},{"name":"cfg.GetMongoTimeout"}]},{"name":"collection.NewProvider","deps":[{"name":"mongoPkg.NewClient","deps":[{"name":"strings.ReplaceAll","deps":[{"name":"cfg.Delay.Mongo.ConnectionString","value":"mongodb://srv_cdp:{{PASSWORD}}@mongo-45-1.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-2.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-3.cdp.dc-2.dev.dcwp.pl:27071/?replicaSet=cdp\u0026readPreference=primary\u0026serverSelectionTimeoutMS=5000\u0026connectTimeoutMS=10000\u0026authSource=admin\u0026authMechanism=SCRAM-SHA-256"},{"name":"cfg.Delay.Mongo.Password"}]},{"name":"cfg.GetMongoTimeout"}]},{"name":"cfg.Delay.Mongo.DB","value":"delay"},{"name":"cfg.Delay.Mongo.RunTimeStateCollection","value":"runtime_state"},{"name":"cfg.GetMongoTimeout"}]},{"name":"time.Now"},{"name":"","deps":[{"name":"\"PROFILE\"","deps":[{"name":"profile.Deserialize"}]},{"name":"\"DICT\"","deps":[{"name":"dict.Deserialize"}]},{"name":"\"EVENT\"","deps":[{"name":"eventReader.Deserialize"}]},{"name":"\"SEGMENT\"","deps":[{"name":"segmentReader.Deserialize","deps":[{"name":"segment.NewSegmentServiceClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]},{"name":"dp.NewDigitalProfileServiceClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]}]}]}]},{"name":"","deps":[{"name":"\"PROFILE\"","deps":[{"name":"profile2.Deserialize","deps":[{"name":"measured.NewQueue","deps":[{"name":"throttle.NewQueue","deps":[{"name":"concurrent.NewQueue","deps":[{"name":"append","deps":[{"name":"rabbit.NewConnection","deps":[{"name":"","deps":[{"name":"Host","deps":[{"name":"strings.Split","deps":[{"name":"cfg.Rabbit.Host","value":"node-1.rabbitmq-1.dc-2.dev.dcwp.pl"}]}]},{"name":"User","deps":[{"name":"cfg.Rabbit.User","value":"cdp"}]},{"name":"Pass","deps":[{"name":"cfg.Rabbit.Password"}]},{"name":"VHost","deps":[{"name":"cfg.Rabbit.Vhost","value":"CDP"}]},{"name":"QueueName","deps":[{"name":"cfg.Rabbit.QueueName","value":"attrChange"}]},{"name":"DlqName","deps":[{"name":"cfg.Rabbit.DlqName","value":"attrChangeDLQ"}]},{"name":"ContentType"}]}]}]}]},{"name":"strings.Split","deps":[{"name":"cfg.AlwaysWriteACIDs","value":"257c138f-ac95-425e-b629-7e28720ba17e,3b75dd46-1dce-46eb-b701-af400da7583d,474283bf-baf0-436e-b7d4-1219c0d22319"}]}]}]}]}]},{"name":"\"CONSOLE\"","deps":[{"name":"consoleWriter.Deserialize"}]}]}]},{"name":"strings.Split","deps":[{"name":"cfg.AlwaysWriteACIDs","value":"257c138f-ac95-425e-b629-7e28720ba17e,3b75dd46-1dce-46eb-b701-af400da7583d,474283bf-baf0-436e-b7d4-1219c0d22319"}]}]},{"name":"cacheScenario.NewRepository","deps":[{"name":"collection.NewProvider","deps":[{"name":"mongoPkg.NewClient","deps":[{"name":"strings.ReplaceAll","deps":[{"name":"cfg.Delay.Mongo.ConnectionString","value":"mongodb://srv_cdp:{{PASSWORD}}@mongo-45-1.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-2.cdp.dc-2.dev.dcwp.pl:27071,mongo-45-3.cdp.dc-2.dev.dcwp.pl:27071/?replicaSet=cdp\u0026readPreference=primary\u0026serverSelectionTimeoutMS=5000\u0026connectTimeoutMS=10000\u0026authSource=admin\u0026authMechanism=SCRAM-SHA-256"},{"name":"cfg.Delay.Mongo.Password"}]},{"name":"cfg.GetMongoTimeout"}]},{"name":"cfg.Delay.Mongo.DB","value":"delay"},{"name":"cfg.Delay.Mongo.RecipeCollection","value":"recipes"},{"name":"cfg.GetMongoTimeout"}]}]}]}]},{"name":"writerfactory.New","deps":[{"name":"","deps":[{"name":"\"CONSOLE\"","deps":[{"name":"consoleWriter.NewWriterFactory"}]},{"name":"\"PROFILE\"","deps":[{"name":"profile2.NewWriterFactory","deps":[{"name":"measured.NewQueue","deps":[{"name":"throttle.NewQueue","deps":[{"name":"concurrent.NewQueue","deps":[{"name":"append","deps":[{"name":"rabbit.NewConnection","deps":[{"name":"","deps":[{"name":"Host","deps":[{"name":"strings.Split","deps":[{"name":"cfg.Rabbit.Host","value":"node-1.rabbitmq-1.dc-2.dev.dcwp.pl"}]}]},{"name":"User","deps":[{"name":"cfg.Rabbit.User","value":"cdp"}]},{"name":"Pass","deps":[{"name":"cfg.Rabbit.Password"}]},{"name":"VHost","deps":[{"name":"cfg.Rabbit.Vhost","value":"CDP"}]},{"name":"QueueName","deps":[{"name":"cfg.Rabbit.QueueName","value":"attrChange"}]},{"name":"DlqName","deps":[{"name":"cfg.Rabbit.DlqName","value":"attrChangeDLQ"}]},{"name":"ContentType"}]}]}]}]},{"name":"strings.Split","deps":[{"name":"cfg.AlwaysWriteACIDs","value":"257c138f-ac95-425e-b629-7e28720ba17e,3b75dd46-1dce-46eb-b701-af400da7583d,474283bf-baf0-436e-b7d4-1219c0d22319"}]}]}]}]}]}]}]}]}]}]},{"name":"cfg.GetCacheCdpIDsTime"},{"name":"cfg.GetCacheCdpIDsCleanupTime"}]}]}]}]}]},{"name":"Closer","deps":[{"name":"append","deps":[{"name":"append","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]},{"name":"measured.NewQueue","deps":[{"name":"throttle.NewQueue","deps":[{"name":"concurrent.NewQueue","deps":[{"name":"append","deps":[{"name":"rabbit.NewConnection","deps":[{"name":"","deps":[{"name":"Host","deps":[{"name":"strings.Split","deps":[{"name":"cfg.Rabbit.Host","value":"node-1.rabbitmq-1.dc-2.dev.dcwp.pl"}]}]},{"name":"User","deps":[{"name":"cfg.Rabbit.User","value":"cdp"}]},{"name":"Pass","deps":[{"name":"cfg.Rabbit.Password"}]},{"name":"VHost","deps":[{"name":"cfg.Rabbit.Vhost","value":"CDP"}]},{"name":"QueueName","deps":[{"name":"cfg.Rabbit.QueueName","value":"attrChange"}]},{"name":"DlqName","deps":[{"name":"cfg.Rabbit.DlqName","value":"attrChangeDLQ"}]},{"name":"ContentType"}]}]}]}]},{"name":"strings.Split","deps":[{"name":"cfg.AlwaysWriteACIDs","value":"257c138f-ac95-425e-b629-7e28720ba17e,3b75dd46-1dce-46eb-b701-af400da7583d,474283bf-baf0-436e-b7d4-1219c0d22319"}]}]}]}]}]},{"name":"ExampleEvent","deps":[{"name":"exampleEventProcessor.Event"}]},{"name":"RecipeRepository","deps":[{"name":"measured.NewRecipeRepo","deps":[{"name":"cache.New","deps":[{"name":"recipe.NewRecipeClient","deps":[{"name":"grpc.Dial","deps":[{"name":"cfg.CoreAPI.Host","value":"core-api.default.svc.cluster.local:9000"},{"name":"grpc.WithTransportCredentials","deps":[{"name":"insecure.NewCredentials"}]},{"name":"grpc.WithUnaryInterceptor","deps":[{"name":"prometheus.UnaryInterceptor"}]}]}]},{"name":"cfg.Meta.Application","value":"px"},{"name":"cfg.GetCacheRecipesTime"},{"name":"cfg.GetCacheRecipesCleanupTime"}]}]}]}]};
let elements = [];

let id = 0;

function processDependency(dep, parentID=null) {
	let myID = id;
	id++;
	elements.push({
		data: {
			id: myID,
			label: dep.name + (dep.value ? ` (${dep.value})` : ''),
			value: dep.value
		}
	});


	if (parentID) {
		elements.push({
			data: {
				source: parentID,
				target: myID,
			}
		});
	}

	if (dep.deps) {
		dep.deps.forEach(d => {
			processDependency(d, myID);
		});
	}
}

processDependency(graphData);

let cy = cytoscape({
	container: document.getElementById('cy'),
	elements: elements,
	style: [
		{
			selector: 'node',
			style: {
				'label': 'data(label)',
				'text-wrap': 'wrap',
				'text-max-width': '60px',
				'background-color': '#666',
				'color': '#fff',
				'font-size': '10px'
			}
		},
		{
			selector: 'edge',
			style: {
				'curve-style': 'bezier',
				'target-arrow-shape': 'triangle'
			}
		}
	],
	layout: {
		name: 'elk',
		nodeDimensionsIncludeLabels: true,
		fit: true,
		padding: 50,
		animate: true,
		elk: {
			algorithm: 'layered'
		}
}});

// Function to get all descendants of a node
function getDescendants(node) {
	let descendants = node.successors(); // This gets all nodes and edges that are reachable from the node
	return descendants;
}

cy.on('tap', 'node', function(evt){
	let node = evt.target;
	let descendants = getDescendants(node);

	if (descendants.some(ele => ele.hidden())) {
		// If any of the descendants are hidden, show them
		descendants.show();
	} else {
		// Otherwise, hide the descendants
		descendants.hide();
	}
});
