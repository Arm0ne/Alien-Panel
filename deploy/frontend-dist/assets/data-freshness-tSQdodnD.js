import{$r as e,Bi as t,Br as n,Ca as r,Fn as i,Gi as a,Hi as o,It as s,Jn as c,Ki as l,Mn as u,Qi as d,Ri as f,Vi as p,Wi as m,Yn as h,_r as g,bn as _,di as v,ea as y,fi as b,ga as x,hn as S,ii as C,li as w,lr as T,mi as E,na as D,nr as O,pi as k,qi as A,ra as j,ua as M,ui as N,va as P,zi as F}from"./router-CGznAFcf.js";import{d as I}from"./index-CjmCa0zB.js";function L(t){let{textColor2:n,primaryColorHover:r,primaryColorPressed:i,primaryColor:a,infoColor:o,successColor:s,warningColor:c,errorColor:l,baseColor:u,borderColor:d,opacityDisabled:f,tagColor:p,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=t;return Object.assign(Object.assign({},I),{closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:f,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:A,textColorCheckable:n,textColorHoverCheckable:n,textColorPressedCheckable:n,textColorChecked:u,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:a,colorCheckedHover:r,colorCheckedPressed:i,border:`1px solid ${d}`,textColor:n,color:p,colorBordered:`rgb(250, 250, 252)`,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${e(a,{alpha:.3})}`,textColorPrimary:a,colorPrimary:e(a,{alpha:.12}),colorBorderedPrimary:e(a,{alpha:.1}),closeIconColorPrimary:a,closeIconColorHoverPrimary:a,closeIconColorPressedPrimary:a,closeColorHoverPrimary:e(a,{alpha:.12}),closeColorPressedPrimary:e(a,{alpha:.18}),borderInfo:`1px solid ${e(o,{alpha:.3})}`,textColorInfo:o,colorInfo:e(o,{alpha:.12}),colorBorderedInfo:e(o,{alpha:.1}),closeIconColorInfo:o,closeIconColorHoverInfo:o,closeIconColorPressedInfo:o,closeColorHoverInfo:e(o,{alpha:.12}),closeColorPressedInfo:e(o,{alpha:.18}),borderSuccess:`1px solid ${e(s,{alpha:.3})}`,textColorSuccess:s,colorSuccess:e(s,{alpha:.12}),colorBorderedSuccess:e(s,{alpha:.1}),closeIconColorSuccess:s,closeIconColorHoverSuccess:s,closeIconColorPressedSuccess:s,closeColorHoverSuccess:e(s,{alpha:.12}),closeColorPressedSuccess:e(s,{alpha:.18}),borderWarning:`1px solid ${e(c,{alpha:.35})}`,textColorWarning:c,colorWarning:e(c,{alpha:.15}),colorBorderedWarning:e(c,{alpha:.12}),closeIconColorWarning:c,closeIconColorHoverWarning:c,closeIconColorPressedWarning:c,closeColorHoverWarning:e(c,{alpha:.12}),closeColorPressedWarning:e(c,{alpha:.18}),borderError:`1px solid ${e(l,{alpha:.23})}`,textColorError:l,colorError:e(l,{alpha:.1}),colorBorderedError:e(l,{alpha:.08}),closeIconColorError:l,closeIconColorHoverError:l,closeIconColorPressedError:l,closeColorHoverError:e(l,{alpha:.12}),closeColorPressedError:e(l,{alpha:.18})})}var R={name:`Tag`,common:S,self:L},z={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},B=N(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[b(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),v(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),v(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),v(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),v(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),b(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[v(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),v(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),b(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),b(`icon, avatar`,[b(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),b(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),b(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[k(`disabled`,[w(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[k(`checked`,`color: var(--n-text-color-hover-checkable);`)]),w(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[k(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),b(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[k(`disabled`,[w(`&:hover`,`background-color: var(--n-color-checked-hover);`),w(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),V=Object.assign(Object.assign(Object.assign({},u.props),z),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),H=n(`n-tag`),U=l({name:`Tag`,props:V,slots:Object,setup(e){let t=x(null),{mergedBorderedRef:n,mergedClsPrefixRef:r,inlineThemeDisabled:a,mergedRtlRef:o,mergedComponentPropsRef:s}=h(e),l=f(()=>e.size||s?.value?.Tag?.size||`medium`),d=u(`Tag`,`-tag`,B,R,e,r);j(H,{roundRef:P(e,`round`)});function p(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function m(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&T(n,t)}}let _={setTextContent(e){let{value:n}=t;n&&(n.textContent=e)}},v=i(`Tag`,o,r),y=f(()=>{let{type:t,color:{color:r,textColor:i}={}}=e,a=l.value,{common:{cubicBezierEaseInOut:o},self:{padding:s,closeMargin:c,borderRadius:u,opacityDisabled:f,textColorCheckable:p,textColorHoverCheckable:m,textColorPressedCheckable:h,textColorChecked:g,colorCheckable:_,colorHoverCheckable:v,colorPressedCheckable:y,colorChecked:b,colorCheckedHover:x,colorCheckedPressed:S,closeBorderRadius:w,fontWeightStrong:T,[E(`colorBordered`,t)]:D,[E(`closeSize`,a)]:O,[E(`closeIconSize`,a)]:k,[E(`fontSize`,a)]:A,[E(`height`,a)]:j,[E(`color`,t)]:M,[E(`textColor`,t)]:N,[E(`border`,t)]:P,[E(`closeIconColor`,t)]:F,[E(`closeIconColorHover`,t)]:I,[E(`closeIconColorPressed`,t)]:L,[E(`closeColorHover`,t)]:R,[E(`closeColorPressed`,t)]:z}}=d.value,B=C(c);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${j} - 8px)`,"--n-bezier":o,"--n-border-radius":u,"--n-border":P,"--n-close-icon-size":k,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":w,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":O,"--n-color":r||(n.value?D:M),"--n-color-checkable":_,"--n-color-checked":b,"--n-color-checked-hover":x,"--n-color-checked-pressed":S,"--n-color-hover-checkable":v,"--n-color-pressed-checkable":y,"--n-font-size":A,"--n-height":j,"--n-opacity-disabled":f,"--n-padding":s,"--n-text-color":i||N,"--n-text-color-checkable":p,"--n-text-color-checked":g,"--n-text-color-hover-checkable":m,"--n-text-color-pressed-checkable":h}}),b=a?c(`tag`,f(()=>{let t=``,{type:r,color:{color:i,textColor:a}={}}=e;return t+=r[0],t+=l.value[0],i&&(t+=`a${g(i)}`),a&&(t+=`b${g(a)}`),n.value&&(t+=`c`),t}),y,e):void 0;return Object.assign(Object.assign({},_),{rtlEnabled:v,mergedClsPrefix:r,contentRef:t,mergedBordered:n,handleClick:p,handleCloseClick:m,cssVars:a?void 0:y,themeClass:b?.themeClass,onRender:b?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:a,onRender:o,$slots:s}=this;o?.();let c=O(s.avatar,e=>e&&A(`div`,{class:`${t}-tag__avatar`},e)),l=O(s.icon,e=>e&&A(`div`,{class:`${t}-tag__icon`},e));return A(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:a,[`${t}-tag--avatar`]:c,[`${t}-tag--icon`]:l,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},l||c,A(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?A(_,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:a,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?A(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),W={key:0,class:`ml-4px`},G={key:0,class:`mt-2px text-12px opacity-80`},K=l({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let n=e,i=x(Date.now()),c,l=f(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),u=f(()=>{if(l.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-l.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,r=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>r?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),h=f(()=>l.value===null?``:new Date(l.value).toLocaleString(`zh-CN`,{hour12:!1}));return y(()=>{n.compact||(c=setInterval(()=>{i.value=Date.now()},6e4))}),d(()=>{c&&clearInterval(c)}),(n,i)=>{let c=U,l=s;return D(),t(l,null,{trigger:M(()=>[a(c,{size:`small`,type:u.value.type},{default:M(()=>[m(r(u.value.label)+` `,1),!e.compact&&h.value?(D(),o(`span`,W,`· `+r(h.value),1)):p(``,!0)]),_:1},8,[`type`])]),default:M(()=>[F(`div`,null,r(u.value.detail),1),h.value?(D(),o(`div`,G,`数据时间：`+r(h.value),1)):p(``,!0)]),_:1})}}});export{U as n,K as t};