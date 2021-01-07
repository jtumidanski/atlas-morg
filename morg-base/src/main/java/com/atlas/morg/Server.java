package com.atlas.morg;

import java.net.URI;

import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.morg.event.consumer.MapCharacterConsumer;
import com.atlas.morg.event.consumer.MonsterDamageConsumer;
import com.atlas.morg.event.consumer.MonsterMovementConsumer;
import com.atlas.morg.processor.MonsterProcessor;
import com.atlas.morg.rest.constant.RestConstants;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.UriBuilder;

public class Server {
   public static void main(String[] args) {
      Runtime.getRuntime().addShutdownHook(new Thread(MonsterProcessor::destroyAll));

      SimpleEventConsumerFactory.create(new MonsterMovementConsumer());
      SimpleEventConsumerFactory.create(new MapCharacterConsumer());
      SimpleEventConsumerFactory.create(new MonsterDamageConsumer());

      URI uri = UriBuilder.host(RestConstants.SERVICE).uri();
      RestServerFactory.create(uri, "com.atlas.morg.rest");
   }
}
