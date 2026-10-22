import{Ai as e,Bn as t,Cn as n,Dt as r,En as i,Fi as a,Gi as o,Ii as s,Ji as c,Li as l,Mi as u,Ni as d,Ri as f,Ui as p,Wr as m,Yr as h,ai as g,cr as _,dn as v,ei as y,er as b,ii as x,ji as S,jr as C,ki as w,la as T,ma as E,ni as D,on as O,qi as k,qn as A,ri as j,sa as M,ta as N,ti as P,zn as F}from"./router-DC4J4Cs2.js";import{d as I}from"./index-CFN0RikG.js";function L(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:i,infoColor:a,successColor:o,warningColor:s,errorColor:c,baseColor:l,borderColor:u,opacityDisabled:d,tagColor:f,closeIconColor:p,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=e;return Object.assign(Object.assign({},I),{closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:d,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:A,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:l,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:i,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${u}`,textColor:t,color:f,colorBordered:`rgb(250, 250, 252)`,closeIconColor:p,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${m(i,{alpha:.3})}`,textColorPrimary:i,colorPrimary:m(i,{alpha:.12}),colorBorderedPrimary:m(i,{alpha:.1}),closeIconColorPrimary:i,closeIconColorHoverPrimary:i,closeIconColorPressedPrimary:i,closeColorHoverPrimary:m(i,{alpha:.12}),closeColorPressedPrimary:m(i,{alpha:.18}),borderInfo:`1px solid ${m(a,{alpha:.3})}`,textColorInfo:a,colorInfo:m(a,{alpha:.12}),colorBorderedInfo:m(a,{alpha:.1}),closeIconColorInfo:a,closeIconColorHoverInfo:a,closeIconColorPressedInfo:a,closeColorHoverInfo:m(a,{alpha:.12}),closeColorPressedInfo:m(a,{alpha:.18}),borderSuccess:`1px solid ${m(o,{alpha:.3})}`,textColorSuccess:o,colorSuccess:m(o,{alpha:.12}),colorBorderedSuccess:m(o,{alpha:.1}),closeIconColorSuccess:o,closeIconColorHoverSuccess:o,closeIconColorPressedSuccess:o,closeColorHoverSuccess:m(o,{alpha:.12}),closeColorPressedSuccess:m(o,{alpha:.18}),borderWarning:`1px solid ${m(s,{alpha:.35})}`,textColorWarning:s,colorWarning:m(s,{alpha:.15}),colorBorderedWarning:m(s,{alpha:.12}),closeIconColorWarning:s,closeIconColorHoverWarning:s,closeIconColorPressedWarning:s,closeColorHoverWarning:m(s,{alpha:.12}),closeColorPressedWarning:m(s,{alpha:.18}),borderError:`1px solid ${m(c,{alpha:.23})}`,textColorError:c,colorError:m(c,{alpha:.1}),colorBorderedError:m(c,{alpha:.08}),closeIconColorError:c,closeIconColorHoverError:c,closeIconColorPressedError:c,closeColorHoverError:m(c,{alpha:.12}),closeColorPressedError:m(c,{alpha:.18})})}var R={name:`Tag`,common:O,self:L},z={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},B=P(`tag`,`
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
`,[j(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),D(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),D(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),D(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),D(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),j(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[D(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),D(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),j(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),j(`icon, avatar`,[j(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),j(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),j(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[x(`disabled`,[y(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[x(`checked`,`color: var(--n-text-color-hover-checkable);`)]),y(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[x(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),j(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[x(`disabled`,[y(`&:hover`,`background-color: var(--n-color-checked-hover);`),y(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),V=Object.assign(Object.assign(Object.assign({},n.props),z),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),H=C(`n-tag`),U=l({name:`Tag`,props:V,slots:Object,setup(e){let r=M(null),{mergedBorderedRef:a,mergedClsPrefixRef:o,inlineThemeDisabled:s,mergedRtlRef:l,mergedComponentPropsRef:u}=t(e),d=w(()=>e.size||u?.value?.Tag?.size||`medium`),f=n(`Tag`,`-tag`,B,R,e,o);c(H,{roundRef:T(e,`round`)});function p(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function m(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&b(n,t)}}let v={setTextContent(e){let{value:t}=r;t&&(t.textContent=e)}},y=i(`Tag`,l,o),x=w(()=>{let{type:t,color:{color:n,textColor:r}={}}=e,i=d.value,{common:{cubicBezierEaseInOut:o},self:{padding:s,closeMargin:c,borderRadius:l,opacityDisabled:u,textColorCheckable:p,textColorHoverCheckable:m,textColorPressedCheckable:_,textColorChecked:v,colorCheckable:y,colorHoverCheckable:b,colorPressedCheckable:x,colorChecked:S,colorCheckedHover:C,colorCheckedPressed:w,closeBorderRadius:T,fontWeightStrong:E,[g(`colorBordered`,t)]:D,[g(`closeSize`,i)]:O,[g(`closeIconSize`,i)]:k,[g(`fontSize`,i)]:A,[g(`height`,i)]:j,[g(`color`,t)]:M,[g(`textColor`,t)]:N,[g(`border`,t)]:P,[g(`closeIconColor`,t)]:F,[g(`closeIconColorHover`,t)]:I,[g(`closeIconColorPressed`,t)]:L,[g(`closeColorHover`,t)]:R,[g(`closeColorPressed`,t)]:z}}=f.value,B=h(c);return{"--n-font-weight-strong":E,"--n-avatar-size-override":`calc(${j} - 8px)`,"--n-bezier":o,"--n-border-radius":l,"--n-border":P,"--n-close-icon-size":k,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":T,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":O,"--n-color":n||(a.value?D:M),"--n-color-checkable":y,"--n-color-checked":S,"--n-color-checked-hover":C,"--n-color-checked-pressed":w,"--n-color-hover-checkable":b,"--n-color-pressed-checkable":x,"--n-font-size":A,"--n-height":j,"--n-opacity-disabled":u,"--n-padding":s,"--n-text-color":r||N,"--n-text-color-checkable":p,"--n-text-color-checked":v,"--n-text-color-hover-checkable":m,"--n-text-color-pressed-checkable":_}}),S=s?F(`tag`,w(()=>{let t=``,{type:n,color:{color:r,textColor:i}={}}=e;return t+=n[0],t+=d.value[0],r&&(t+=`a${_(r)}`),i&&(t+=`b${_(i)}`),a.value&&(t+=`c`),t}),x,e):void 0;return Object.assign(Object.assign({},v),{rtlEnabled:y,mergedClsPrefix:o,contentRef:r,mergedBordered:a,handleClick:p,handleCloseClick:m,cssVars:s?void 0:x,themeClass:S?.themeClass,onRender:S?.onRender})},render(){var e;let{mergedClsPrefix:t,rtlEnabled:n,closable:r,color:{borderColor:i}={},round:a,onRender:o,$slots:s}=this;o?.();let c=A(s.avatar,e=>e&&f(`div`,{class:`${t}-tag__avatar`},e)),l=A(s.icon,e=>e&&f(`div`,{class:`${t}-tag__icon`},e));return f(`div`,{class:[`${t}-tag`,this.themeClass,{[`${t}-tag--rtl`]:n,[`${t}-tag--strong`]:this.strong,[`${t}-tag--disabled`]:this.disabled,[`${t}-tag--checkable`]:this.checkable,[`${t}-tag--checked`]:this.checkable&&this.checked,[`${t}-tag--round`]:a,[`${t}-tag--avatar`]:c,[`${t}-tag--icon`]:l,[`${t}-tag--closable`]:r}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},l||c,f(`span`,{class:`${t}-tag__content`,ref:`contentRef`},(e=this.$slots).default?.call(e)),!this.checkable&&r?f(v,{clsPrefix:t,class:`${t}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:a,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?f(`div`,{class:`${t}-tag__border`,style:{borderColor:i}}):null)}}),W={key:0,class:`ml-4px`},G={key:0,class:`mt-2px text-12px opacity-80`},K=l({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(t){let n=t,i=M(Date.now()),c,l=w(()=>{if(!n.dataAt)return null;let e=new Date(n.dataAt).getTime();return Number.isFinite(e)?e:null}),f=w(()=>{if(l.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-l.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(n.delayedAfterMinutes,1)*60*1e3,r=Math.max(n.expiredAfterMinutes,n.delayedAfterMinutes)*60*1e3;return e>r?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(r/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),m=w(()=>l.value===null?``:new Date(l.value).toLocaleString(`zh-CN`,{hour12:!1}));return o(()=>{n.compact||(c=setInterval(()=>{i.value=Date.now()},6e4))}),p(()=>{c&&clearInterval(c)}),(n,i)=>{let o=U,c=r;return k(),S(c,null,{trigger:N(()=>[s(o,{size:`small`,type:f.value.type},{default:N(()=>[a(E(f.value.label)+` `,1),!t.compact&&m.value?(k(),d(`span`,W,`· `+E(m.value),1)):u(``,!0)]),_:1},8,[`type`])]),default:N(()=>[e(`div`,null,E(f.value.detail),1),m.value?(k(),d(`div`,G,`数据时间：`+E(m.value),1)):u(``,!0)]),_:1})}}});export{U as n,K as t};