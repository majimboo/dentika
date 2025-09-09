<template>
  <Transition 
    :name="transitionName"
    :mode="mode"
    :appear="appear"
    @before-enter="onBeforeEnter"
    @enter="onEnter"
    @after-enter="onAfterEnter"
    @before-leave="onBeforeLeave"
    @leave="onLeave"
    @after-leave="onAfterLeave"
  >
    <slot />
  </Transition>
</template>

<script>
export default {
  name: 'BaseTransition',
  props: {
    name: {
      type: String,
      default: 'fade'
    },
    mode: {
      type: String,
      default: '',
      validator: (value) => ['', 'out-in', 'in-out'].includes(value)
    },
    appear: {
      type: Boolean,
      default: false
    }
  },
  emits: [
    'before-enter',
    'enter',
    'after-enter',
    'before-leave',
    'leave',
    'after-leave'
  ],
  computed: {
    transitionName() {
      return this.name
    }
  },
  methods: {
    onBeforeEnter(el) {
      this.$emit('before-enter', el)
    },
    onEnter(el) {
      this.$emit('enter', el)
    },
    onAfterEnter(el) {
      this.$emit('after-enter', el)
    },
    onBeforeLeave(el) {
      this.$emit('before-leave', el)
    },
    onLeave(el) {
      this.$emit('leave', el)
    },
    onAfterLeave(el) {
      this.$emit('after-leave', el)
    }
  }
}
</script>

<style>
/* Fade Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Slide Down Transition */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* Slide Up Transition */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* Scale Transition */
.scale-enter-active,
.scale-leave-active {
  transition: all 0.2s ease;
}

.scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* Bounce Transition */
.bounce-enter-active {
  animation: bounce-in 0.5s;
}

.bounce-leave-active {
  animation: bounce-in 0.5s reverse;
}

@keyframes bounce-in {
  0% {
    transform: scale(0.3);
    opacity: 0;
  }
  50% {
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.9);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

/* List Transition */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-leave-active {
  position: absolute;
}
</style>