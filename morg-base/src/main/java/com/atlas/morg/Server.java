package com.atlas.morg;

import com.atlas.kafka.consumer.SimpleEventConsumerFactory;
import com.atlas.morg.event.consumer.MapCharacterConsumer;
import com.atlas.morg.event.consumer.MonsterDamageConsumer;
import com.atlas.morg.event.consumer.MonsterMovementConsumer;
import com.atlas.shared.rest.RestServerFactory;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import java.net.URI;

public class Server {
   public static void main(String[] args) {
      SimpleEventConsumerFactory.create(new MonsterMovementConsumer());
      SimpleEventConsumerFactory.create(new MapCharacterConsumer());
      SimpleEventConsumerFactory.create(new MonsterDamageConsumer());

      URI uri = UriBuilder.host(RestService.MONSTER_REGISTRY).uri();
      RestServerFactory.create(uri, "com.atlas.morg.rest");
   }
}
