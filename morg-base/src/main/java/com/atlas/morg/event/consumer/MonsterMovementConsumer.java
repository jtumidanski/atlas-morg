package com.atlas.morg.event.consumer;

import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.MonsterMovementEvent;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.MonsterRegistry;

public class MonsterMovementConsumer implements SimpleEventHandler<MonsterMovementEvent> {
   @Override
   public void handle(Long aLong, MonsterMovementEvent event) {
      MonsterRegistry.getInstance().moveMonster(event.uniqueId(), event.endX(), event.endY(), event.stance());
   }

   @Override
   public Class<MonsterMovementEvent> getEventClass() {
      return MonsterMovementEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Monster Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_MONSTER_MOVEMENT);
   }
}
