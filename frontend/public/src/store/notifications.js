import {FeedService} from "../service";
import {findIndexById} from "../util";

const NOTIFICATION_LIMIT = 5;

export const NotificationsStore = {
    state: {
        notifications: [],
        notificationCount: 0
    },
    mutations: {
        setNotifications(state, {notifications, count}) {
            state.notifications = notifications;
            state.notificationCount = count;
        },
        setNotification(state, {index, notification}) {
            state.notifications[index] = notification;
        }
    },
    actions: {
        loadNotifications(context) {
            FeedService.getNotifications(false, NOTIFICATION_LIMIT)
                .then(notifications => {
                    if(notifications.count > 0) {
                        let lastNotification = window.localStorage.getItem("last_notification");
                        let lastNotificationDate = new Date(lastNotification);
                        let notificationDate = new Date(notifications.feed[0].def_time);

                        if(!lastNotification || lastNotificationDate < notificationDate) {
                            window.localStorage.setItem("last_notification", notifications.feed[0].def_time);
                            context.dispatch("showBrowserNotification", notifications.feed[0]);
                        }
                    }

                    context.commit("setNotifications", {notifications: notifications.feed, count: notifications.count});
                })
                .catch(e => {
                    console.error(e);
                });
        },
        showBrowserNotification(context, notification) {
            if(!Notification) {
                return;
            }

            if(Notification.permission === "granted") {
                let dom = document.createElement("div");
                dom.innerHTML = notification.notification;
                let text = dom.textContent || dom.innerText || "";
                let title = `${notification.triggered_by_user.firstname} ${notification.triggered_by_user.lastname} ${text}`;
                let notify = new Notification(title, {
                    icon: "/static/img/favicon-32x32.png"
                });
                notify.onclick = () => {
                    window.focus();
                };
            }
            else if(Notification.permission !== "denied") {
                Notification.requestPermission();
            }
        },
        toggleNotificationRead(context, id) {
            let index = findIndexById(context.state.notifications, id);

            if(index > -1) {
                let notification = context.state.notifications[index];
                notification.read = !notification.read;
                context.commit("setNotification", {index, notification});
                let count = context.state.notificationCount;
                count = notification.read ? count-1 : count+1;
                context.commit("setNotifications", {notifications: context.state.notifications, count});
            }
        },
        markNotificationsRead(context) {
            let notifications = context.state.notifications;

            for(let i = 0; i < notifications.length; i++) {
                notifications[i].read = true;
            }

            context.commit("setNotifications", {notifications, count: 0});
        }
    },
    getters: {
        notifications(state) {
            return state.notifications;
        },
        notificationCount(state) {
            return state.notificationCount;
        }
    }
};
